package versatilis

import (
	"crypto/rand"
	"errors"
	"log"
	"net"
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
}

// creates a new TCP connection.  if initiator is true, then it Dials the
// address specified by dst.  if initiator is false, it takes in a net.Conn object associated
// with an already existing TCP connection.
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
	c := &Conn{
		directionality:       ChannelTypeBidirectional,
		addressType:          AddressTypeTCP,
		listenAddress:        &addr,
		dstAddress:           &addr,
		transportLayerInBuf:  NewVBuffer(),
		inBuf:                NewVBuffer(),
		inChan:               make(chan bool),
		transportLayerOutBuf: NewVBuffer(),
		outBuf:               NewVBuffer(),
		outChan:              make(chan bool),
	}
	go c.outBufProcessor()
	return c, nil
}

// sends some data via an already-established connection.  or, more precisely,
// kicks off the process of the data being transported according to some
// transport protocol
func (conn *Conn) Send(b []byte) (n int, err error) {
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
		plaintext, err := conn.outBuf.ReadAll()
		if err != nil {
			log.Fatalf("cannot read buffer: %v", err)
		}

		// encrypt the plaintext
		ad := make([]byte, 16)
		rand.Read(ad)
		ciphertext, err := state.encryptState.Encrypt(nil, ad, plaintext)
		if err != nil {
			return -1, err
		}

		// TODO: create two ProtoBuf objects, a PackageHdr type with the length
		// of the latter, where the latter contains the ciphertext
	}
}
