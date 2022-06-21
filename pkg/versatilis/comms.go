package versatilis

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"

	"google.golang.org/protobuf/proto"
)

type VersatilisDialer struct{}

func send(dst *Address, p *Package) error {
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

func recv(listenAddress *Address, block bool) (*Package, error) {
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
		conn, ok := listenAddress.EndPoint.(*net.Conn)
		if !ok {
			return nil, errors.New("invalid address")
		}
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

func (state *State) Write(b any) (n int, err error) {
	buf := MessageBuffer{}
	m := &Message{
		Id:      "write",
		Payload: b,
	}
	buf = append(buf, m)
	if err := state.Send(state.dst, &buf); err != nil {
		return 0, err
	}
	return 0, nil
}
