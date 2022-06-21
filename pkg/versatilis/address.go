package versatilis

import "net"

type AddressType int64

const (
	AddressTypeInvalid AddressType = iota
	AddressTypeTCP
	AddressTypeSCTP
	AddressTypeDHT
	AddressTypeChan
)

// implements the net.Addr interface (see https://pkg.go.dev/net#Addr)
type Address struct {
	Type     AddressType
	EndPoint any
}

func (a Address) Network() string {
	switch a.Type {
	case AddressTypeTCP:
		return "versatilis-TCP"
	case AddressTypeSCTP:
		return "versatilis-SCTP"
	case AddressTypeDHT:
		return "versatilis-DHT"
	case AddressTypeChan:
		return "versatilis-GoChan"
	default:
		return "ERROR"
	}
}

func (a Address) String() string {
	switch a.Type {
	case AddressTypeTCP:
		c, ok := a.EndPoint.(*net.Conn)
		if !ok {
			return "ERROR "
		}
		return (*c).RemoteAddr().String()
	default:
		return "ERROR"
	}
}
