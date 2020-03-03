package main

// Credit to @mlowicki.
//
// https://medium.com/t/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c

import (
	"code.cloudfoundry.org/bytefmt"
	"context"
	"crypto/tls"
	"github.com/juju/ratelimit"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"time"
)

func handleTunneling(w http.ResponseWriter, r *http.Request, bucket *ratelimit.Bucket) {
	conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)

	destConn := RateLimitedConn{conn, bucket}

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func copyWithLog(destination io.Writer, source io.Reader, description string) {
	start := time.Now()

	if bytes, err := io.Copy(destination, source); err != nil {
		log.Warn(err)
	} else {
		size := bytefmt.ByteSize(uint64(bytes))
		duration := time.Since(start)
		rate := bytefmt.ByteSize(uint64(float64(bytes)/duration.Seconds())) + "/s"

		log.WithFields(log.Fields{"Size": size, "Duration": duration, "Rate": rate}).Info(description)
	}
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()

	copyWithLog(destination, source, "Transfered")
}

func handleHTTP(w http.ResponseWriter, req *http.Request, bucket *ratelimit.Bucket) {
	dialer := &net.Dialer{}

	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := dialer.DialContext(ctx, network, addr)

		return RateLimitedConn{conn, bucket}, err
	}

	resp, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	copyWithLog(w, resp.Body, req.RequestURI)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func proxy(port string, rate int, done <-chan bool) {
	log.Info("Goforward listening on :" + port + " with ratelimit " + bytefmt.ByteSize(uint64(rate)))

	bucket := ratelimit.NewBucketWithRate(float64(rate), int64(rate))

	server := &http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.WithFields(log.Fields{"Method": r.Method, "RemoteAddr": r.RemoteAddr}).Info(r.RequestURI)

			if r.Method == http.MethodConnect {
				handleTunneling(w, r, bucket)
			} else {
				handleHTTP(w, r, bucket)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	go server.ListenAndServe()

	<-done

	log.Info("Exiting Goforward")

	server.Close()
}
