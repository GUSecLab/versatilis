package versatilis

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// creates a new TCP connection.  if initiator is true, then it Dials the
// address specified by dst.  if initiator is false, it takes in a net.Conn
// object associated with an already existing TCP connection.
func NewTCPConn(initiator bool, existingConn *net.Conn, dst string) (*Conn, error) {
	addr := Address{
		Type: AddressTypeTCP,
	}
	switch initiator {
	case true:
		conn, err := net.Dial("tcp", dst)
		if err != nil {
			return nil, err
		}
		addr.EndPoint = &conn

	case false:
		addr.EndPoint = existingConn
	}

	// create a new connection
	c, err := newConn(initiator)
	if err != nil {
		return nil, err
	}

	// set some TCP-specific fields
	c.directionality = ChannelTypeBidirectional
	c.addressType = AddressTypeTCP
	c.listenAddress = &addr
	c.dstAddress = &addr

	// kick off the necessary buffer processors
	go c.outBufProcessor()
	go c.inBufProcessor()
	go tcpTransportSend(c)
	go tcpTransportRecv(c)

	// let's do a handshake!
	c.handleHandshake()

	return c, nil
}

// a helper function that sends data whenever there's data to be sent
func tcpTransportSend(conn *Conn) {
	for range conn.outChanTransport {
		// we have data to send!
		buf, err := conn.transportOutBuf.ReadAll()
		if err != nil {
			log.Errorf("[%v] send error: %v", conn, err)
			continue
		}

		sock, ok := conn.dstAddress.EndPoint.(*net.Conn)
		if !ok {
			log.Errorf("[%v] send error", conn)
			continue
		}

		// send data until there's no data left to send
		for len(buf) > 0 {
			n, err := (*sock).Write(buf) // send it over the wire
			if err != nil {
				log.Errorf("[%v] send error", conn)
				break
			}
			buf = buf[n:]
		}
	}
}

// a helper function that sends a message to the inChanTransport channel
// whenever data has been received
func tcpTransportRecv(conn *Conn) {

	const bufSize = 65535

	sock, ok := conn.dstAddress.EndPoint.(*net.Conn)
	if !ok {
		log.Errorf("[%v] recv error", conn)
		return
	}

	buf := make([]byte, bufSize)

	// listen indefinitely
	for {
		n, err := (*sock).Read(buf)
		if err != nil {
			log.Errorf("[%v] read error: %v", conn, err)
			return
		}
		if n > 0 {
			bytesToWrite := buf[:n]
			log.Infof("[%v] [recv] %v", conn, bytesToWrite)
			bytesWritten, err := conn.transportInBuf.Write(bytesToWrite)
			if err != nil {
				log.Errorf("[%v] buffer write error: %v", conn, err)
				return
			}
			if bytesWritten != n {
				log.Errorf("[%v] buffer write error", conn)
				return
			}

			log.Infof("[%v] wrote %v bytes to transportInBuf", conn, bytesWritten)

			// notify something that there's data available to be read
			conn.inChanTransport <- true
		}
	}

}
