package versatilis

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"net"

	"google.golang.org/protobuf/proto"
)

type VersatilisDialer struct{}

type readBuffer []byte

// internal helper function that actually sends a Package via the appropriate
// transport protocol to the address specified by Address
func sendPackageViaTransport(dst *Address, p *Package) error {
	switch dst.Type {
	case AddressTypeChan:
		ch, ok := dst.EndPoint.(chan *Package)
		if !ok {
			return errors.New("invalid address")
		}
		ch <- p // send it!
		return nil
	case AddressTypeTCP:
		conn, ok := dst.EndPoint.(*net.Conn)
		if !ok {
			return errors.New("invalid address")
		}
		buf, err := proto.Marshal(p)
		if err != nil {
			return err
		}
		l := make([]byte, 2)
		binary.LittleEndian.PutUint16(l, uint16(len(buf)))
		(*conn).Write(l)
		(*conn).Write(buf)
		return nil
	default:
		return errors.New("unsupport address type")
	}
}

// internal helper function that actually receives a Package via the appropriate
// transport protocol.  The function listens for incoming traffic from the
// specified listenAddress
func recvPackageViaTransport(listenAddress *Address, block bool) (*Package, error) {
	switch listenAddress.Type {
	case AddressTypeChan:
		ch, ok := listenAddress.EndPoint.(chan *Package)
		if !ok {
			return nil, errors.New("invalid address")
		}
		if block {
			p := <-ch
			return p, nil
		} else {
			// we don't block, so do a select
			select {
			case p := <-ch:
				return p, nil // data is available
			default:
				return nil, nil // no data available
			}
		}

	case AddressTypeTCP:
		if !block {
			return nil, errors.New("non-blocking mode not supported for TCP")
		}
		conn, ok := listenAddress.EndPoint.(*net.Conn)
		if !ok {
			return nil, errors.New("invalid address")
		}

		// TODO: SUPPORT NON-BLOCKING MODE!!!

		lbe := make([]byte, 2)
		if n, err := io.ReadFull(*conn, lbe); err != nil || n != 2 {
			return nil, errors.New("invalid message")
		}
		var l uint16
		buf := bytes.NewReader(lbe)
		if err := binary.Read(buf, binary.LittleEndian, &l); err != nil {
			return nil, err
		}
		bigbuf := make([]byte, l)
		if n, err := io.ReadFull(*conn, bigbuf); err != nil || n != int(l) {
			return nil, errors.New("invalid message")
		}
		p := &Package{}
		if err := proto.Unmarshal(bigbuf, p); err != nil {
			return nil, err
		}
		return p, nil

	default:
		return nil, errors.New("unsupport address type")
	}
}

// receives a message or returns nil, nil if no message is available (if block is false)
func (state *State) Receive(block bool) (*MessageBuffer, error) {

	p, err := recvPackageViaTransport(state.listenAddress, block)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	plaintext, err := state.decryptState.Decrypt(nil, p.NoiseAuthTag, p.NoiseCiphertext)
	if err != nil {
		return nil, err
	}
	var messages MessageBuffer
	err = proto.Unmarshal(plaintext, &messages)
	if err != nil {
		return nil, err
	} else {
		return &messages, nil
	}
}

// write arbitrary byte stream via Versatilis channel
func (state *State) Write(b []byte) (n int, err error) {

	// encapsulate bytes into a MessageBuffer of Messages
	buffer := MessageBuffer{}
	m := &Message{
		Payload: &Message_BytesMessage{b},
	}
	buffer.Messages = append(buffer.Messages, m)

	// the plaintext is the serialized buffer
	plaintext, err := proto.Marshal(&buffer)
	if err != nil {
		return -1, err
	}

	// encrypt the plaintext
	ad := make([]byte, 16)
	rand.Read(ad)
	ciphertext, err := state.encryptState.Encrypt(nil, ad, plaintext)
	if err != nil {
		return -1, err
	}

	// send the ciphertext via the transport
	if err := sendPackageViaTransport(state.dst, &Package{
		Version:         Version.String(),
		NoiseCiphertext: ciphertext,
		NoiseAuthTag:    ad,
	}); err != nil {
		return -1, err
	}

	return len(b), nil
}

// read arbitrary byte stream from Versatilis channel
// TODO: this function really needs improvement
func (state *State) Read(b []byte) (n int, err error) {
	n = 0

	// handle the simple case
	if len(b) == 0 {
		return
	}

	// first, see if there's data in the read buffer that we can read
	state.readMu.Lock()
	if len(state.readBuffer) > 0 {
		bytesToCopy := minInt(len(state.readBuffer), len(b))
		copyInPlace(b, state.readBuffer, bytesToCopy)
		n += bytesToCopy
		state.readBuffer = state.readBuffer[bytesToCopy:]
		// TODO: not optimal  :(
		state.readMu.Unlock()
		return
	}
	state.readMu.Unlock()

	// if there isn't data in the buffer already, get data over the network
	messages, err := state.Receive(true)
	if err != nil {
		return -1, err
	}

	// next, copy the data into the readBuffer
	state.readMu.Lock()
	defer state.readMu.Unlock()
	for _, message := range messages.Messages {
		switch x := message.Payload.(type) {
		case *Message_BytesMessage:
			state.readBuffer = append(state.readBuffer, x.BytesMessage...)
		case *Message_StringMessage:
			barr := []byte(x.StringMessage)
			state.readBuffer = append(state.readBuffer, barr...)
		case nil: // field not set
			// do nothing
		default:
			return -1, errors.New("invalid message received")
		}
	}

	// finally, copy bytes out of the buffer
	if len(state.readBuffer) > 0 {
		bytesToCopy := minInt(len(state.readBuffer), len(b))
		copyInPlace(b, state.readBuffer, bytesToCopy)
		state.readBuffer = state.readBuffer[bytesToCopy:]
		n += bytesToCopy
		return
	}
	return
}

/*
func (state *State) Dial(t AddressType, address any) (net.Conn, error) {
	if address == nil {
		return nil, errors.New("invalid address")
	}
	switch t {
	case AddressTypeTCP:
		s, ok := address.(string)
		if !ok {
			return nil, errors.New("invalid address")
		}
		state.connected = true
		return net.Dial("tcp", s)
	default:
		return nil, errors.New("unsupported address type")
	}
}
*/
