package versatilis

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"

	"google.golang.org/protobuf/proto"

	"github.com/flynn/noise"
)

var constPackageHdrSize int
var constHandshakeHdrSize int

type ChannelType uint64

const (
	ChannelTypeUndefined ChannelType = iota
	ChannelTypeBidirectional
	ChannelTypeInbound
	ChannelTypeOutbound
)

type Conn struct {

	// pointer to its parent Versatilis instance
	v *Versatilis

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
	transportInBuf *vBuffer

	// (authenticated) plaintext data ready to be read by an app
	inBuf *vBuffer

	// a channel used to signal that data is available for reading
	inChan chan bool

	// a channel used to signal that data is available for reading via the transport
	inChanTransport chan bool

	// encrypted and marshalled data ready to be sent off the wire
	transportOutBuf *vBuffer

	// plaintext data ready to be encrypted and marshalled and moved over to
	// transportLayerOutBuf
	outBuf *vBuffer

	// a channel used to signal that data is available in the outBuf for processing
	outChan chan bool

	// a channel used to signal that data is available for sending via the transport
	outChanTransport chan bool

	noiseConfig         *noise.Config
	noiseHandshakeState *noise.HandshakeState

	// if true, the handshake is completed
	handshakeCompleted bool

	// a channel for conveying that the handshake is done
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

func (conn *Conn) String() string {
	switch conn.directionality {
	case ChannelTypeBidirectional, ChannelTypeInbound:
		return fmt.Sprintf("[conn (type %v) (initiator: %v)]",
			conn.listenAddress, conn.noiseConfig.Initiator)
	case ChannelTypeOutbound:
		return fmt.Sprintf("[conn (type %v) (initiator: %v)]",
			conn.dstAddress, conn.noiseConfig.Initiator)
	default:
		return "[conn (invalid/unititialized)]"
	}
}

// a helper function which sets some initial/default values for a new Conn, and
// does some init stuff for Noise
func newConn(v *Versatilis, initiator bool) (conn *Conn, err error) {

	conn = &Conn{
		v:                      v,
		transportInBuf:         NewVBuffer(),
		inBuf:                  NewVBuffer(),
		inChan:                 make(chan bool),
		inChanTransport:        make(chan bool),
		transportOutBuf:        NewVBuffer(),
		outBuf:                 NewVBuffer(),
		outChan:                make(chan bool),
		outChanTransport:       make(chan bool),
		handshakeCompleted:     false,
		badState:               false,
		badStateErr:            nil,
		handshakeCompletedChan: make(chan bool),
	}

	constPackageHdrSize = proto.Size(&PackageHdr{PackageSize: 1})
	constHandshakeHdrSize = proto.Size(&HandshakeMsgHdr{Size: 1})

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

// returns true iff there is data available for reading.  This data could
// correspond to a handshake message or a payload.
func (conn *Conn) IsDataAvailable() bool {
	return (conn.inBuf.Size() > 0)
}

// sends some data via an already-established connection.  or, more precisely,
// kicks off the process of the data being transported according to some
// transport protocol
func (conn *Conn) Send(b []byte) (n int, err error) {
	if conn.badState {
		return -1, conn.badStateErr
	}

	if conn == nil || conn.directionality == ChannelTypeUndefined || conn.directionality == ChannelTypeInbound {
		return -1, errors.New("invalid channel")
	}
	n, err = conn.outBuf.Write(b)
	if err != nil {
		return -1, err
	}
	conn.outChan <- true // signal that more data is available
	return
}

// Reads some plaintext data via an already-established connection.
func (conn *Conn) Read(maxBytes int) ([]byte, error) {
	if conn.badState {
		return nil, conn.badStateErr
	}

	if conn == nil || conn.directionality == ChannelTypeUndefined || conn.directionality == ChannelTypeOutbound {
		return nil, errors.New("invalid channel")
	}

	// block and wait for data to arrive
	for conn.inBuf.Size() <= 0 {
		<-conn.inChan
	}

	b, err := conn.inBuf.Read(maxBytes)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// a helper function which writes either a package (if p is not null) or a
// handshake message (if h is not nil) to the transportOutBuf
func (conn *Conn) writePackageToTransportOutBuf(p *Package, h *HandshakeMsg) (err error) {
	if p != nil && h != nil {
		panic("invalid usage")
	}

	var marshalledData []byte
	var marshalledDataLen []byte

	switch {
	case p != nil:
		marshalledData, err = proto.Marshal(p)
		if (err != nil) || len(marshalledData) >= math.MaxUint32 {
			return errors.New("data handling error")
		}
		hdr := PackageHdr{
			PackageSize: uint32(len(marshalledData)),
		}
		marshalledDataLen, err = proto.Marshal(&hdr)
		if (err != nil) || len(marshalledDataLen) != constPackageHdrSize {
			return errors.New("data handling error")
		}

	case h != nil:
		marshalledData, err = proto.Marshal(h)
		if (err != nil) || len(marshalledData) >= math.MaxUint32 {
			return errors.New("data handling error")
		}
		hdr := HandshakeMsgHdr{
			Size: uint32(len(marshalledData)),
		}
		marshalledDataLen, err = proto.Marshal(&hdr)
		if (err != nil) || len(marshalledDataLen) != constHandshakeHdrSize {
			return errors.New("data handling error")
		}
	}

	if n, err := conn.transportOutBuf.Write(marshalledDataLen); err != nil || n != len(marshalledDataLen) {
		return errors.New("data handling error")
	}
	if n, err := conn.transportOutBuf.Write(marshalledData); err != nil || n != len(marshalledData) {
		return errors.New("data handling error")
	}

	return nil
}

// this helper function waits for a signal that data is available, and then
// packagizes the data and pushes the packages to the transportLayerOutBuf
func (conn *Conn) outBufProcessor() {

	if !conn.handshakeCompleted {
		<-conn.handshakeCompletedChan
	}

	// wait for a signal that data is available
	for range conn.outChan {

		plaintext, err := conn.outBuf.ReadAll()
		if err != nil {
			conn.v.Log.Fatalf("cannot read buffer: %v", err)
		}
		if len(plaintext) == 0 {
			continue
		}

		// if we get here, there's actual plaintext, so let's encrypt it with
		// Noise's AEAD scheme

		var p Package
		p.Authtag = make([]byte, 16)
		rand.Read(p.Authtag)
		p.Ciphertext, err = conn.encryptState.Encrypt(nil, p.Authtag, plaintext)
		if err != nil {
			conn.badState = true
			conn.badStateErr = err
			return // don't process any more messages!
		}

		// write the package to the transport out buf
		if err := conn.writePackageToTransportOutBuf(&p, nil); err != nil {
			conn.badState = true
			conn.badStateErr = err
			return
		}

		// lastly, inform some send function (e.g., tcpTransportSend) that it
		// should look at the transportOutBuf
		conn.outChanTransport <- true

	}
}

// this helper function "wakes up" when data is available in the transportInBuf,
// and then de-packages and decrypts the data and puts it in the inBuf
func (conn *Conn) inBufProcessor() {

	if !conn.handshakeCompleted {
		<-conn.handshakeCompletedChan
	}

	// wait for a signal that data is available
	for range conn.inChanTransport {

		// first, we read the package header from transportInBuf
		if conn.transportInBuf.Size() < int(constPackageHdrSize) {
			// there's not enough data to read a full PackageHdr
			continue
		}
		phRaw, err := conn.transportInBuf.Read(int(constPackageHdrSize))
		if err != nil {
			conn.v.Log.Errorf("buffer read error: %v", err)
			continue
		}
		var ph PackageHdr
		if err := proto.Unmarshal(phRaw, &ph); err != nil {
			conn.v.Log.Errorf("unmarshalling error: %v", err)
			continue
		}

		// let's wait until we have enough bytes to read
		for conn.transportInBuf.Size() < int(ph.PackageSize) {
			conn.v.Log.Infof("[%v] inBufProcessor waiting for signal on conn.inChanTransport", conn)
			<-conn.inChanTransport
		}

		// woo-hoo!  we have enough bytes to read out package
		pRaw, err := conn.transportInBuf.Read(int(ph.PackageSize))
		if err != nil {
			conn.v.Log.Errorf("buffer read error: %v", err)
			continue
		}
		var p Package
		if err := proto.Unmarshal(pRaw, &p); err != nil {
			conn.v.Log.Errorf("unmarshalling error: %v", err)
			continue
		}

		// now, let's decrypt the package
		plaintext, err := conn.decryptState.Decrypt(nil, p.Authtag, p.Ciphertext)
		if err != nil {
			conn.v.Log.Errorf("decryption error: %v", err)
			continue
		}
		conn.v.Log.Debugf("decrypted: %v", plaintext)

		// finally, let's write the plaintext to the inbuf and signal that data
		// is available
		if _, err := conn.inBuf.Write(plaintext); err != nil {
			conn.v.Log.Errorf("cannot write to inbuf: %v", err)
			continue
		}
		conn.inChan <- true

	}
}

// this helper function blocks until it can return a handshake message from the
// transport
func (conn *Conn) getHandshakeMsgFromTransport() (*HandshakeMsg, error) {

	var hh HandshakeMsgHdr
	var h HandshakeMsg

	conn.v.Log.Debugf("[%v] getHandshakeMsgFromTransport: getting handshake message from transport", conn)

	// first, we wait for enough data at the inbound transport buffer to grab a
	// HandshakeMsgHdr
	for conn.transportInBuf.Size() < int(constHandshakeHdrSize) {
		conn.v.Log.Debugf("[%v] not enuf bytes for handshake message; have %v, want %v",
			conn, conn.transportInBuf.Size(), constHandshakeHdrSize)
		<-conn.inChanTransport
	}

	// grab the handshakemsghdr
	hhRaw, err := conn.transportInBuf.Read(int(constHandshakeHdrSize))
	if err != nil {
		conn.v.Log.Errorf("buffer read error: %v", err)
		return nil, err
	}
	if err := proto.Unmarshal(hhRaw, &hh); err != nil {
		conn.v.Log.Errorf("unmarshalling error: %v", err)
		return nil, err
	}

	// let's wait until we have enough bytes to read the full handshakeMsg
	for conn.transportInBuf.Size() < int(hh.Size) {
		conn.v.Log.Debugf("[%v] getHandshakeMsgFromTransport waiting for signal on conn.inChanTransport", conn)
		<-conn.inChanTransport
	}

	// woo-hoo!  we have enough bytes to read in our handshake message
	hRaw, err := conn.transportInBuf.Read(int(hh.Size))
	if err != nil {
		conn.v.Log.Errorf("buffer read error: %v", err)
		return nil, err
	}
	if err := proto.Unmarshal(hRaw, &h); err != nil {
		conn.v.Log.Errorf("unmarshalling error: %v", err)
		return nil, err
	}

	conn.v.Log.Debugf("[%v] getHandshakeMsgFromTransport: got handshake message from transport!", conn)

	return &h, nil
}

func (conn *Conn) handleHandshake() {
	var handshakeMsgRaw []byte
	var err error

	out := make([]byte, 0, 4096)

	for _, action := range conn.handshakeSendRecvPattern {
		if action { // send event

			conn.v.Log.Debugf("[%v] prepping to send handshake message", conn)

			handshakeMsgRaw, conn.encryptState, conn.decryptState, err =
				conn.noiseHandshakeState.WriteMessage(out, nil)
			if err != nil {
				conn.v.Log.Errorf("handshake error [send]: %v", err)
				conn.badState = true
				conn.badStateErr = err
				return
			}

			h := HandshakeMsg{
				Message: handshakeMsgRaw,
			}
			// write the package to the transport out buf
			if err := conn.writePackageToTransportOutBuf(nil, &h); err != nil {
				conn.badState = true
				conn.badStateErr = err
				return
			}
			// lastly, inform some send function (e.g., tcpTransportSend) that it
			// should look at the transportOutBuf
			conn.outChanTransport <- true

			conn.v.Log.Debugf("[%v] sent handshake message", conn)
		} else { // receive event

			conn.v.Log.Debugf("[%v] prepping to receive handshake message", conn)

			handshakeMsg, err := conn.getHandshakeMsgFromTransport()
			if err != nil {
				conn.badState = true
				conn.badStateErr = err
				return
			}

			_, conn.decryptState, conn.encryptState, err =
				conn.noiseHandshakeState.ReadMessage(out, handshakeMsg.Message)
			if err != nil {
				conn.v.Log.Errorf("handshake error [recv]: %v", err)
				conn.badState = true
				conn.badStateErr = err
				return
			}

			conn.v.Log.Debugf("[%v] received handshake message", conn)
		}
	}

	if err != nil {
		panic(err)
	}
	if conn.encryptState != nil {
		conn.handshakeCompleted = true
		conn.handshakeCompletedChan <- true
		// if it's bidirectional, we need to info both the inbufprocessor and
		// the outbufprocessor
		if conn.directionality == ChannelTypeBidirectional {
			conn.handshakeCompletedChan <- true
		}
		conn.v.Log.Infof("[%v] handshake complete", conn)
	}
}
