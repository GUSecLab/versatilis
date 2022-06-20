package versatilis

import "errors"

func send(dst *Address, p *Package) error {
	switch dst.Type {
	case AddressTypeChan:
		ch, ok := dst.EndPoint.(chan *Package)
		if !ok {
			return errors.New("invalid address")
		}
		ch <- p // send it!
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
	default:
		return nil, errors.New("unsupport address type")
	}
}
