package main

import (
	"net"
	"time"
	versatilis "versatilis/pkg/versatilis"

	log "github.com/sirupsen/logrus"
)

func initiator(dst *versatilis.Address, listenAddress *versatilis.Address) {
	v := versatilis.New(true, "initiator")
	log.Infof("I am %v", v.Name)

	v.DoHandshake(dst, listenAddress)

	var n int
	var err error
	b := make([]byte, 50)

	if n, err = v.Read(b); err != nil {
		panic(err)
	}
	b = b[:n]
	log.Infof("initiator received %v", string(b))

	/*
		mb, err := v.Receive(listenAddress, true)
		if err != nil {
			panic(err)
		} else {
			log.Info("initiator received: %v", mb)
		}
	*/
}

func initiatorChan(done chan bool, toInitiator chan *versatilis.Package, toResponder chan *versatilis.Package) {
	log.Debug("starting initiator")

	listenAddress := &versatilis.Address{
		Type:     versatilis.AddressTypeChan,
		EndPoint: toInitiator,
	}
	dst := &versatilis.Address{
		Type:     versatilis.AddressTypeChan,
		EndPoint: toResponder,
	}

	initiator(dst, listenAddress)

	done <- true
}

func initiatorTCP(done chan bool) {
	log.Debug("[tcp] starting TCP initiator")

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		panic(err)
	} else {
		defer conn.Close()
	}

	addressConn := &versatilis.Address{
		Type:     versatilis.AddressTypeTCP,
		EndPoint: &conn,
	}

	initiator(addressConn, addressConn)

	done <- true
}

func responder(dst *versatilis.Address, listenAddress *versatilis.Address, msg string) {
	v := versatilis.New(false, "responder")
	log.Infof("I am %v", v.Name)
	v.DoHandshake(dst, listenAddress)

	v.Write([]byte(msg))
}

func responderChan(done chan bool, toInitiator chan *versatilis.Package, toResponder chan *versatilis.Package) {
	log.Debug("starting responder")

	listenAddress := &versatilis.Address{
		Type:     versatilis.AddressTypeChan,
		EndPoint: toResponder,
	}
	dst := &versatilis.Address{
		Type:     versatilis.AddressTypeChan,
		EndPoint: toInitiator,
	}

	responder(dst, listenAddress, "testing")

	done <- true
}

func responderTCP(done chan bool) {
	log.Debug("[tcp] starting TCP responder")

	ln, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		panic(err)
	}

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	} else {
		defer conn.Close()
		defer ln.Close()
	}
	log.Debugf("[tcp] connection established from %v", conn)

	addressConn := &versatilis.Address{
		Type:     versatilis.AddressTypeTCP,
		EndPoint: &conn,
	}

	responder(addressConn, addressConn, "hello world")

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

	log.Info("waiting 1 second for next test")
	time.Sleep(time.Second * 1)

	go responderTCP(done)
	time.Sleep(time.Second * 1)
	go initiatorTCP(done)

	<-done
	<-done

}
