package listener

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/dynamicgo/mesh"
	"github.com/dynamicgo/mesh/libp2p/netwrapper"
	"github.com/dynamicgo/slf4go"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	protocol "github.com/libp2p/go-libp2p-protocol"
)

type libp2pListener struct {
	sync.Mutex
	slf4go.Logger
	addr     string
	conn     chan net.Conn
	ctx      context.Context
	cancel   context.CancelFunc
	protocol protocol.ID
	host     host.Host
}

// Listen listen protocol stream incoming
func Listen(ctx context.Context, host host.Host, protocol protocol.ID) net.Listener {
	ctx, cancel := context.WithCancel(ctx)

	listener := &libp2pListener{
		Logger:   slf4go.Get("libp2p-listener"),
		conn:     make(chan net.Conn),
		ctx:      ctx,
		cancel:   cancel,
		protocol: protocol,
		host:     host,
	}

	host.SetStreamHandler(listener.protocol, listener.streamHandler)

	return listener
}

func (listener *libp2pListener) streamHandler(stream inet.Stream) {
	conn := &netwrapper.StreamConn{Stream: stream}

	for {
		if listener.sendStream(conn) {
			return
		}

		time.Sleep(time.Second)
	}
}

func (listener *libp2pListener) sendStream(conn net.Conn) bool {
	listener.Lock()
	defer listener.Unlock()

	select {
	case <-listener.ctx.Done():
		listener.WarnF("handle input stream on closed listener %s", listener.addr)
		return true
	default:
	}

	select {
	case listener.conn <- conn:
		return true
	default:
		return false
	}
}

func (listener *libp2pListener) Accept() (net.Conn, error) {
	conn, ok := <-listener.conn

	if !ok {
		return nil, mesh.ErrNetworkClosed
	}

	return conn, nil
}

func (listener *libp2pListener) Close() error {
	listener.Lock()
	defer listener.Unlock()

	listener.host.RemoveStreamHandler(listener.protocol)

	listener.cancel()
	close(listener.conn)

	return nil
}

func (listener *libp2pListener) Addr() net.Addr {
	return &netwrapper.NetAddr{
		Addr: listener.addr,
	}
}
