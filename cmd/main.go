package main

import (
	"github.com/guseclab/versatilis/pkg/versatilis"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("starting up")

	versatilis.New()

}
