package main

import (
	"time"
	versatilis "versatilis/pkg/versatilis"

	log "github.com/sirupsen/logrus"
)

func initiator(done chan bool, toInitiator chan *versatilis.Package, toResponder chan *versatilis.Package) {
	log.Debug("starting initiator")
	v := versatilis.New(true, "initiator")
	log.Infof("I am %v", v.Name)
	v.DoHandshake(toResponder, toInitiator)

	// ...
	for i := 0; i < 15; i++ {
		m, err := v.Receive(toInitiator, false)
		if err != nil {
			panic(err)
		} else {
			if m != nil {
				log.Infof("initiator received %v", m)
			} else {
				log.Info("Initiator received no message")
			}
		}
		time.Sleep(time.Millisecond * 100)
	}

	done <- true
}

func responder(done chan bool, toInitiator chan *versatilis.Package, toResponder chan *versatilis.Package) {
	log.Debug("starting responder")
	v := versatilis.New(false, "responder")
	log.Infof("I am %v", v.Name)
	v.DoHandshake(toInitiator, toResponder)
	done <- true

	// ...
	buf := versatilis.MessageBuffer{}
	for x := 0; x < 10; x++ {
		m := &versatilis.Message{
			Id: "type1",
			Payload: struct {
				First   string
				Last    string
				Counter int
			}{
				First:   "Micah",
				Last:    "Sherr",
				Counter: x,
			},
		}
		buf = append(buf, m)
	}

	if err := v.Send(toInitiator, &buf); err != nil {
		panic(err)
	}

}

func main() {

	log.SetLevel(log.DebugLevel)
	versatilis.SetLogLevel(log.InfoLevel)

	log.Info("starting up")

	done := make(chan bool)
	toInitiator := make(chan *versatilis.Package)
	toResponder := make(chan *versatilis.Package)

	go initiator(done, toInitiator, toResponder)
	go responder(done, toInitiator, toResponder)

	<-done
	<-done
}
