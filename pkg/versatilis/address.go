package versatilis

type AddressType int64

const (
	AddressTypeInvalid AddressType = iota
	AddressTypeTCP
	AddressTypeSCTP
	AddressTypeIP
	AddressTypeHostname
	AddressTypeDHT
	AddressTypeChan
)

type Address struct {
	Type     AddressType
	EndPoint any
}
