package main

// https://medium.com/t/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c

import (
	"crypto/tls"
    "flag"
    "io"
    "net"
    "net/http"
	"time"
	"github.com/juju/ratelimit"
	log "github.com/sirupsen/logrus"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
    destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
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

func transfer(destination io.WriteCloser, source io.ReadCloser) {
    defer destination.Close()
	defer source.Close()

	start := time.Now()
	
	bucket := ratelimit.NewBucketWithRate(100*1024, 100*1024)
	if bytes, err := io.Copy(destination, ratelimit.Reader(source, bucket)); err != nil {
		log.Error(err)
	} else {
		log.WithFields(log.Fields{"Bytes": bytes, "Duration": time.Since(start)}).Info("Copied")
	}
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
    resp, err := http.DefaultTransport.RoundTrip(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }
    defer resp.Body.Close()
    copyHeader(w.Header(), resp.Header)
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
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

	log.Info("Goforward Proxy Listening on :" + port)
	
	server := &http.Server{
        Addr: ":" + port,
        Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.WithFields(log.Fields{"Method": r.Method, "RemoteAddr": r.RemoteAddr}).Info(r.RequestURI)

            if r.Method == http.MethodConnect {
                handleTunneling(w, r)
            } else {
                handleHTTP(w, r)
            }
        }),
        // Disable HTTP/2.
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	
    if proto == "http" {
        log.Fatal(server.ListenAndServe())
	}
}