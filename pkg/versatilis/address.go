package versatilis

type AddressType int64

const (
	AddressTypeInvalid AddressType = iota
	AddressTypeNetConn
	AddressTypeIP
	AddressTypeHostname
	AddressTypeDHT
	AddressTypeChan
)

type Address struct {
	Type     AddressType
	EndPoint interface{}
}
