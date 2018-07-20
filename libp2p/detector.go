package libp2p

import (
	"strings"

	"github.com/dynamicgo/slf4go"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	multiaddr "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multistream"
)

type protocolDetector struct {
	slf4go.Logger
	host     host.Host
	detector ProtocolDetector
	name     string
}

// ProtocolDetector .
type ProtocolDetector interface {
	Name() string
	Protocols() []protocol.ID
	Update(protocol []protocol.ID, id peer.ID)
	Removed(id peer.ID)
}

func (node *meshNode) Detect(detector ProtocolDetector) {
	pd := &protocolDetector{
		host:     node.host,
		detector: detector,
		Logger:   node.Logger,
		name:     node.name(),
	}

	node.DebugF("[%s] create detector %s", node.name(), detector.Name())

	node.host.Network().Notify(pd)
}

func (detector *protocolDetector) protocols() (ps []string) {

	for _, p := range detector.detector.Protocols() {
		ps = append(ps, string(p))
	}

	return
}

func (detector *protocolDetector) Connected(network inet.Network, conn inet.Conn) {

	remotePeer := conn.RemotePeer()
	peerStore := detector.host.Peerstore()

	detector.DebugF("[%s] connected to %s ", detector.name, remotePeer.Pretty())

	selected, err := peerStore.SupportsProtocols(remotePeer, detector.protocols()...)

	if err != nil {
		detector.ErrorF("[%s] query peer %s support protocol err  ", detector.name, remotePeer.Pretty(), err)
	}

	if len(selected) > 0 {

		var protocols []protocol.ID

		for _, p := range selected {
			protocols = append(protocols, protocol.ID(p))
		}

		if network.Connectedness(remotePeer) == inet.Connected {

			detector.DebugF(
				"[%s] update peer %s support protocols [%s]",
				detector.detector.Name(),
				remotePeer,
				strings.Join(selected, ","),
			)

			detector.detector.Update(protocols, remotePeer)
		}

		return
	}

	detector.WarnF("[%s] query peer %s support protocol -- not found", detector.name, remotePeer.Pretty())

	go detector.testProtocols(conn)
}

func (detector *protocolDetector) testProtocols(conn inet.Conn) {

	remotePeer := conn.RemotePeer()
	stream, err := conn.NewStream()

	detector.DebugF("[%s] test peer %s support protocols", detector.name, remotePeer.Pretty())

	if err != nil {
		detector.ErrorF("test peer %s support protocols err %s", remotePeer, err)
		return
	}

	defer stream.Close()

	var selected []string

	for _, p := range detector.protocols() {
		if err := multistream.SelectProtoOrFail(p, stream); err != nil {
			continue
		}

		selected = append(selected, p)
	}

	detector.DebugF("[%s] test peer %s support protocols finish, found %d", detector.name, remotePeer.Pretty(), len(selected))

	detector.host.Peerstore().AddProtocols(remotePeer, selected...)

	if len(selected) > 0 {

		var protocols []protocol.ID

		for _, p := range selected {
			protocols = append(protocols, protocol.ID(p))
		}

		detector.DebugF(
			"[%s] update peer %s support protocols [%s]",
			detector.detector.Name(),
			remotePeer,
			strings.Join(selected, ","),
		)

		detector.detector.Update(protocols, remotePeer)
	}

}

func (detector *protocolDetector) Disconnected(network inet.Network, conn inet.Conn) {
	remotePeer := conn.RemotePeer()

	if network.Connectedness(remotePeer) == inet.Connected {
		// We're still connected.
		return
	}

	detector.detector.Removed(remotePeer)
}

func (detector *protocolDetector) Listen(inet.Network, multiaddr.Multiaddr)      {}
func (detector *protocolDetector) ListenClose(inet.Network, multiaddr.Multiaddr) {}
func (detector *protocolDetector) OpenedStream(inet.Network, inet.Stream)        {}
func (detector *protocolDetector) ClosedStream(inet.Network, inet.Stream)        {}
