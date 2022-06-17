package main

import (
	versatilis "versatilis/pkg/versatilis"

	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.DebugLevel)

	log.Info("starting up")

	versatilis.New()

}
