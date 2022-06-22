package versatilis

type ChannelType uint64

const (
	ChannelTypeUndefined ChannelType = iota
	ChannelTypeBidirectional
	ChannelTypeInbound
	ChannelTypeOutbound
)

type Conn struct {
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
	transportLayerInBuf *vBuffer

	// (authenticated) plaintext data ready to be read by an app
	inBuf *vBuffer

	// encrypted and marshalled data ready to be sent off the wire
	transportLayerOutBuf *vBuffer

	// plaintext data ready to be encrypted and marshalled and moved over to
	// transportLayerOutBuf
	outBuf *vBuffer
}

func NewConn(directionality ChannelType, addressType AddressType) *Conn {
	return &Conn{
		directionality:       directionality,
		addressType:          addressType,
		listenAddress:        nil,
		dstAddress:           nil,
		transportLayerInBuf:  NewVBuffer(),
		inBuf:                NewVBuffer(),
		transportLayerOutBuf: NewVBuffer(),
		outBuf:               NewVBuffer(),
	}
}
