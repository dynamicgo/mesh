package netwrapper

import (
	"net"

	inet "github.com/libp2p/go-libp2p-net"
	manet "github.com/multiformats/go-multiaddr-net"
)

// NetAddr .
type NetAddr struct {
	Addr string
}

// Network .
func (addr *NetAddr) Network() string {
	return "mesh"
}

// String .
func (addr *NetAddr) String() string {
	return addr.Addr
}

// StreamConn libp2p wrapper to adapter net.Conn
type StreamConn struct {
	inet.Stream
}

// fakeLocalAddr returns a dummy local address.
func fakeLocalAddr() net.Addr {
	localIP := net.ParseIP("127.0.0.1")
	return &net.TCPAddr{IP: localIP, Port: 0}
}

// fakeRemoteAddr returns a dummy remote address.
func fakeRemoteAddr() net.Addr {
	remoteIP := net.ParseIP("127.1.0.1")
	return &net.TCPAddr{IP: remoteIP, Port: 0}
}

// LocalAddr returns the local address.
func (c *StreamConn) LocalAddr() net.Addr {
	addr, err := manet.ToNetAddr(c.Stream.Conn().LocalMultiaddr())
	if err != nil {
		return fakeLocalAddr()
	}
	return addr
}

// RemoteAddr returns the remote address.
func (c *StreamConn) RemoteAddr() net.Addr {
	addr, err := manet.ToNetAddr(c.Stream.Conn().RemoteMultiaddr())
	if err != nil {
		return fakeRemoteAddr()
	}
	return addr
}
