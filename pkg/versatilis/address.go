package versatilis

type AddressType int64

const (
	AddressTypeInvalid AddressType = iota
	AddressTypeNetTCP
	AddressTypeIP
	AddressTypeHostname
	AddressTypeDHT
	AddressTypeChan
)

type Address struct {
	Type     AddressType
	EndPoint interface{}
}
