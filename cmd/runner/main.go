package main

import (
	versatilis "versatilis/pkg/versatilis"

	log "github.com/sirupsen/logrus"
)

func initiator(done chan bool) {
	log.Debug("starting initiator")
	v := versatilis.New(true, "initiator")
	v.GenKey(versatilis.KeyTypeStatic)

	done <- true
}

func responder(done chan bool) {
	log.Debug("starting responder")
	v := versatilis.New(false, "responder")
	v.GenKey(versatilis.KeyTypeStatic)

	done <- true
}

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("starting up")

	done := make(chan bool)

	go initiator(done)
	go responder(done)

	<-done
	<-done
}
