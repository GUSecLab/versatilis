package versatilis

import log "github.com/sirupsen/logrus"

type State struct {
}

func init() {
	log.Debug("initializing Versātilis")
}

func New() *State {
	return new(State)
}
