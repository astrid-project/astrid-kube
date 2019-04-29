package main

import (
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func main() {
	l := log.New()
	l.Infoln("Starting...")

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		l.Infoln("Received an interrupt, stopping services...")
		//cleanup(services, c)
		close(cleanupDone)
	}()
	<-cleanupDone
}
