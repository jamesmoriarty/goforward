package main

import (
	"flag"
	"os"
	"os/signal"
    "syscall"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	var port string
	flag.StringVar(&port, "port", "8888", "Proxy listen port")

	var rate int
	flag.IntVar(&rate, "rate", 512*1024, "Proxy bandwidth ratelimit")

	flag.Parse()

	done := make(chan bool, 1)
	go proxy(port, rate, done)

	// Block until signal and exit
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	done <- true
}
