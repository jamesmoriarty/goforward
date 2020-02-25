package main

// https://medium.com/t/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c

import (
	"code.cloudfoundry.org/bytefmt"
	"context"
	"crypto/tls"
	"flag"
	"github.com/juju/ratelimit"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/http"
	"time"
)

type RateLimitedConn struct {
	net.Conn
	*ratelimit.Bucket
}

func (wrap RateLimitedConn) Read(b []byte) (n int, err error) {
	// start := time.Now()
	wrap.Bucket.Wait(int64(len(b)))
	// duration := time.Since(start)
	// log.WithFields(log.Fields{"Size": len(b), "Duration": duration}).Info("Read")
	return wrap.Conn.Read(b)
}

func (wrap RateLimitedConn) Write(b []byte) (n int, err error) {
	// start := time.Now()
	wrap.Bucket.Wait(int64(len(b)))
	// duration := time.Since(start)
	// log.WithFields(log.Fields{"Size": len(b), "Duration": duration}).Info("Write")
	return wrap.Conn.Write(b)
}

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
		log.Error(err)
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

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	var proto string
	flag.StringVar(&proto, "proto", "http", "Proxy protocol")

	var port string
	flag.StringVar(&port, "port", "8888", "Proxy listen port")

	flag.Parse()

	if proto != "http" {
		log.Fatal("Protocol must be either http")
	}

	log.Info("Goforward listening on :" + port)

	// Bucket adding 512KB every second, holding max 10MB
	bucket := ratelimit.NewBucketWithRate(512*1024, 1024*1024*10)

	server := &http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.WithFields(log.Fields{"Method": r.Method, "RemoteAddr": r.RemoteAddr}).Warn(r.RequestURI)

			if r.Method == http.MethodConnect {
				handleTunneling(w, r, bucket)
			} else {
				handleHTTP(w, r, bucket)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	if proto == "http" {
		log.Fatal(server.ListenAndServe())
	}
}
