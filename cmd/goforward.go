package main

import (
	"flag"
	"github.com/jamesmoriarty/goforward"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	var port string
	flag.StringVar(&port, "port", "8888", "Proxy listen port")
	var rate int
	flag.IntVar(&rate, "rate", 512*1024, "Proxy bandwidth ratelimit")

	flag.Parse()

	shutdown := make(chan bool, 1)

	go goforward.Listen(port, rate, shutdown)

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)

	<-wait

	shutdown <- true
}
