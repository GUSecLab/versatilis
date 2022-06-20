package versatilis

import (
	"errors"
	"io"
	"net"

	"google.golang.org/protobuf/proto"
)

func send(dst *Address, p *Package) error {
	switch dst.Type {
	case AddressTypeChan:
		ch, ok := dst.EndPoint.(chan *Package)
		if !ok {
			return errors.New("invalid address")
		}
		ch <- p // send it!
		return nil
	case AddressTypeNetConn:
		conn, ok := dst.EndPoint.(*net.Conn)
		if !ok {
			return errors.New("invalid address")
		}
		buf, err := proto.Marshal(p)
		if err != nil {
			return err
		}
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

	case AddressTypeNetConn:
		conn, ok := listenAddress.EndPoint.(*net.Conn)
		if !ok {
			return nil, errors.New("invalid address")
		}
		bigbuf := make([]byte, 4096)
		smallbuf := make([]byte, 4096)
		for {
			n, err := (*conn).Read(smallbuf)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return nil, err
				}
			}
			bigbuf = append(bigbuf, smallbuf[:n]...)
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
