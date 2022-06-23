package versatilis

import (
	"crypto/rand"
	"errors"
	"log"
	"math"
	"net"

	"github.com/flynn/noise"
	"google.golang.org/protobuf/proto"
)

type ChannelType uint64

const (
	ChannelTypeUndefined ChannelType = iota
	ChannelTypeBidirectional
	ChannelTypeInbound
	ChannelTypeOutbound
)

type Conn struct {
	// directionality of the channel
	directionality ChannelType

	// if the channel is a ChannelTypeBidirectional or ChannelTypeInbound, the
	// address to listen on.  if the channel is ChannelTypeOutbound, than nil
	listenAddress *Address

	// if the channel is a ChannelTypeBidirectional or ChannelTypeOutbound, the
	// address of our remote peer.  If the channel is ChannelTypeInbound, than nil
	dstAddress *Address

	// specifies the type of this connection (if bidirectional, the type must be
	// the same in both directions).  this is for convenience, and should match
	// the type of listenAddress or dstAddress (whichever is not nil)
	addressType AddressType

	// the raw bytes off the wire (encrypted and marshalled)
	transportLayerInBuf *vBuffer

	// (authenticated) plaintext data ready to be read by an app
	inBuf *vBuffer

	// a channel used to signal that data is available for sending
	inChan chan bool

	// encrypted and marshalled data ready to be sent off the wire
	transportLayerOutBuf *vBuffer

	// plaintext data ready to be encrypted and marshalled and moved over to
	// transportLayerOutBuf
	outBuf *vBuffer

	// a channel used to signal that data is available for sending
	outChan chan bool

	noiseConfig         *noise.Config
	noiseHandshakeState *noise.HandshakeState

	// if true, the handshake is completed
	handshakeCompleted bool

	// a channel that's used for conveying that the handshake process is done
	handshakeCompletedChan chan bool

	// "true"s indicate send events; "false" are receives
	handshakeSendRecvPattern []bool

	encryptState *noise.CipherState
	decryptState *noise.CipherState

	// if true, this connection should not be used
	badState bool

	// and the reason it shouldn't be used
	badStateErr error
}

// a helper function which sets some initial/default values for a new Conn, and
// does some init stuff for Noise
func newConn(initiator bool) (conn *Conn, err error) {

	conn = &Conn{
		transportLayerInBuf:    NewVBuffer(),
		inBuf:                  NewVBuffer(),
		inChan:                 make(chan bool),
		transportLayerOutBuf:   NewVBuffer(),
		outBuf:                 NewVBuffer(),
		outChan:                make(chan bool),
		handshakeCompleted:     false,
		handshakeCompletedChan: make(chan bool),
		badState:               false,
		badStateErr:            nil,
	}

	conn.noiseConfig = &noise.Config{
		CipherSuite: noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256),
		Random:      rand.Reader,
		Initiator:   initiator,
		Pattern:     noise.HandshakeNN, // TODO: this is a bad choice
	}

	// TODO: this is specific to noise.HandshakeNN, which is a bad choice
	if initiator {
		z := [...]bool{true, false}
		conn.handshakeSendRecvPattern = z[:]
	} else {
		z := [...]bool{false, true}
		conn.handshakeSendRecvPattern = z[:]
	}

	if conn.noiseHandshakeState, err = noise.NewHandshakeState(*conn.noiseConfig); err != nil {
		return nil, err
	}
	return conn, nil
}

// creates a new TCP connection.  if initiator is true, then it Dials the
// address specified by dst.  if initiator is false, it takes in a net.Conn object associated
// with an already existing TCP connection.
// TODO: move this function to another file (maybe tcp.go?)
func NewTCPConn(initiator bool, existingConn *net.Conn, dst string) (*Conn, error) {
	addr := Address{
		Type: AddressTypeTCP,
	}
	switch initiator {
	case true:
		conn, err := net.Dial("tcp", dst)
		if err != nil {
			return nil, err
		}
		addr.EndPoint = &conn

	case false:
		addr.EndPoint = existingConn
	}

	// create a new connection
	c, err := newConn(initiator)
	if err != nil {
		return nil, err
	}

	// set some TCP-specific fields
	c.directionality = ChannelTypeBidirectional
	c.addressType = AddressTypeTCP
	c.listenAddress = &addr
	c.dstAddress = &addr

	go c.outBufProcessor()
	return c, nil
}

// sends some data via an already-established connection.  or, more precisely,
// kicks off the process of the data being transported according to some
// transport protocol
func (conn *Conn) Send(b []byte) (n int, err error) {
	if conn.badState {
		return -1, conn.badStateErr
	}

	if conn == nil || conn.directionality == ChannelTypeUndefined {
		return -1, errors.New("invalid channel")
	}
	n, err = conn.outBuf.Write(b)
	if err != nil {
		conn.outChan <- true
	}
	return
}

// this helper function waits for a signal that data is available, and then
// packagizes the data and pushes the packages to the transportLayerOutBuf
func (conn *Conn) outBufProcessor() {
	// wait for a signal that data is available
	for range conn.inChan {

		if !conn.handshakeCompleted {
			// a special case!  we need to complete the handshake.  so we should
			// kickoff a process for that, and utilize a timeout to handle
			// failed handshakes

			// TODO
			// go conn.handleHandshake()

			// TODO: remember to send something on this channel :)
			// wait until the handshake completes
			<-conn.handshakeCompletedChan

			// free some mem?
			close(conn.handshakeCompletedChan)
		}

		plaintext, err := conn.outBuf.ReadAll()
		if err != nil {
			log.Fatalf("cannot read buffer: %v", err)
		}
		if plaintext == nil || len(plaintext) == 0 {
			continue
		}

		// if we get here, there's actual plaintext, so let's encrypt it with
		// Noise's AEAD scheme

		var p Package
		var hdr PackageHdr
		p.Authtag = make([]byte, 16)
		rand.Read(p.Authtag)
		p.Ciphertext, err = conn.encryptState.Encrypt(nil, p.Authtag, plaintext)
		if err != nil {
			conn.badState = true
			conn.badStateErr = err
			return // don't process any more messages!
		}
		marshalledPackage, err := proto.Marshal(&p)
		if (err != nil) || len(marshalledPackage) >= math.MaxUint32 {
			conn.badState = true
			conn.badStateErr = err
			return // don't process any more messages!
		}
		hdr.PackageSize = uint32(len(marshalledPackage))
		marshalledPackageHdr, err := proto.Marshal(&hdr)
		if (err != nil) || len(marshalledPackageHdr) != CONSTANTTODO {
			conn.badState = true
			conn.badStateErr = err
			return // don't process any more messages!
		}

		if n, err := conn.transportLayerOutBuf.Write(marshalledPackageHdr); err != nil || n != len(marshalledPackageHdr) {
			conn.badState = true
			conn.badStateErr = err
			return // don't process any more messages!
		}
		if n, err := conn.transportLayerOutBuf.Write(marshalledPackage); err != nil || n != len(marshalledPackage) {
			conn.badState = true
			conn.badStateErr = err
			return // don't process any more messages!
		}
	}
}
