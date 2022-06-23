package main

import (
	"fmt"
	"net"
	"time"
	versatilis "versatilis/pkg/versatilis"
)

var v *versatilis.Versatilis

func initiatorTCP(done chan bool) {
	v.Log.Debug("[tcp] starting TCP initiator")

	vConn, err := versatilis.NewTCPConn(v, true, nil, "localhost:9999")
	if err != nil {
		panic(err)
	}
	v.Log.Info("launched Versatilis TCP connection")

	for {
		buf, err := vConn.Read(1024)
		if err != nil {
			panic(err)
		}
		v.Log.Infof("MAIN.GO: Received %v", string(buf))
	}

}

func responderTCP(done chan bool) {
	v.Log.Debug("[tcp] starting TCP responder")

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
	v.Log.Debugf("[tcp] connection established from %v", conn)

	vConn, err := versatilis.NewTCPConn(v, false, &conn, "")
	if err != nil {
		panic(err)
	}
	v.Log.Info("launched Versatilis TCP connection")

	for i := 0; i < 5; i++ {
		s := fmt.Sprintf("Hello world; counter is %v", i)
		vConn.Send([]byte(s))
		time.Sleep(time.Second * 2)
	}

	done <- true
}

func main() {

	v = versatilis.New()

	done := make(chan bool)

	v.Log.Info("starting up")

	/*
		v.Log.Info("doing some channel tests")
		toInitiator := make(chan *versatilis.Package)
		toResponder := make(chan *versatilis.Package)

		go initiatorChan(done, toInitiator, toResponder)
		go responderChan(done, toInitiator, toResponder)

		<-done
		<-done
	*/

	go responderTCP(done)
	time.Sleep(time.Second * 1)
	go initiatorTCP(done)

	<-done
	<-done

}
