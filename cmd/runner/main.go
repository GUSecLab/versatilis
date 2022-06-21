package main

import (
	"net"
	"time"
	versatilis "versatilis/pkg/versatilis"

	log "github.com/sirupsen/logrus"
)

func initiatorChan(done chan bool, toInitiator chan *versatilis.Package, toResponder chan *versatilis.Package) {
	log.Debug("starting initiator")
	v := versatilis.New(true, "initiator")
	log.Infof("I am %v", v.Name)

	listenAddress := &versatilis.Address{
		Type:     versatilis.AddressTypeChan,
		EndPoint: toInitiator,
	}
	dst := &versatilis.Address{
		Type:     versatilis.AddressTypeChan,
		EndPoint: toResponder,
	}
	v.DoHandshake(dst, listenAddress)

	incomingSrc := &versatilis.Address{
		Type:     versatilis.AddressTypeChan,
		EndPoint: toInitiator,
	}

	// ...
	for i := 0; i < 15; i++ {
		m, err := v.Receive(incomingSrc, false)
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

func responderChan(done chan bool, toInitiator chan *versatilis.Package, toResponder chan *versatilis.Package) {
	log.Debug("starting responder")
	v := versatilis.New(false, "responder")
	log.Infof("I am %v", v.Name)

	listenAddress := &versatilis.Address{
		Type:     versatilis.AddressTypeChan,
		EndPoint: toResponder,
	}
	dst := &versatilis.Address{
		Type:     versatilis.AddressTypeChan,
		EndPoint: toInitiator,
	}

	v.DoHandshake(dst, listenAddress)

	// ...
	buf := versatilis.MessageBuffer{}
	for x := 0; x < 10; x++ {
		m := &versatilis.Message{
			Id:      "type1",
			Payload: x,
		}
		buf = append(buf, m)
	}

	if err := v.Send(dst, &buf); err != nil {
		panic(err)
	}
	done <- true
}

func initiatorTCP(done chan bool) {
	log.Debug("[tcp] starting TCP initiator")
	v := versatilis.New(true, "initiator")
	log.Infof("[tcp] I am %v", v.Name)

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		panic(err)
	}

	addressConn := &versatilis.Address{
		Type:     versatilis.AddressTypeTCP,
		EndPoint: &conn,
	}

	v.DoHandshake(addressConn, addressConn)

	// ...
	m, err := v.Receive(addressConn, true)
	if err != nil {
		panic(err)
	} else {
		if m != nil {
			log.Infof("[tcp] initiator received %v", m)
		} else {
			log.Info("[tcp] initiator received no message")
		}
	}
	done <- true
}

func responderTCP(done chan bool) {
	log.Debug("[tcp] starting TCP responder")
	v := versatilis.New(false, "responder")
	log.Infof("[tcp] I am %v", v.Name)

	ln, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		panic(err)
	}

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	log.Debugf("[tcp] connection established from %v", conn)

	addressConn := &versatilis.Address{
		Type:     versatilis.AddressTypeTCP,
		EndPoint: &conn,
	}

	v.DoHandshake(addressConn, addressConn)

	// ...
	buf := versatilis.MessageBuffer{}
	m := &versatilis.Message{
		Id:      "type2",
		Payload: "this is a message sent via TCP",
	}
	buf = append(buf, m)

	if err := v.Send(addressConn, &buf); err != nil {
		panic(err)
	}
	done <- true
}

func main() {

	log.SetLevel(log.InfoLevel)
	versatilis.SetLogLevel(log.InfoLevel)

	done := make(chan bool)

	log.Info("starting up")
	log.Info("doing some channel tests")
	toInitiator := make(chan *versatilis.Package)
	toResponder := make(chan *versatilis.Package)

	go initiatorChan(done, toInitiator, toResponder)
	go responderChan(done, toInitiator, toResponder)

	<-done
	<-done

	log.Info("waiting 3 seconds for next test")
	time.Sleep(time.Second * 3)

	go responderTCP(done)
	time.Sleep(time.Second * 5)
	go initiatorTCP(done)

	<-done
	<-done

}
