package versatilis

type AddressType int64

const (
	AddressTypeIP AddressType = iota
	AddressTypeHostname
	AddressTypeDHT
)

type Address struct {
	Type     AddressType
	EndPoint interface{}
}
