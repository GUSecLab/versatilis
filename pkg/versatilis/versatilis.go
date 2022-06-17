package versatilis

import log "github.com/sirupsen/logrus"

type State struct {
}

type Message struct {
	Id      string
	Payload interface{}
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.Debug("initializing Versātilis")
}

func New() *State {
	return new(State)
}
