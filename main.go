package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	var proto string
	flag.StringVar(&proto, "proto", "http", "Proxy protocol")

	var port string
	flag.StringVar(&port, "port", "8888", "Proxy listen port")

	var rate int
	flag.IntVar(&rate, "rate", 512*1024, "Proxy bandwidth ratelimit")

	flag.Parse()

	if proto != "http" {
		log.Fatal("Protocol must be http")
	}

	proxy(proto, port, rate)
}
