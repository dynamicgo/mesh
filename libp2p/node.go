package libp2p

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/dynamicgo/mesh/libp2p/listener"
	"github.com/dynamicgo/mesh/libp2p/netwrapper"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/mesh"
	"github.com/dynamicgo/mesh/addr"
	"github.com/dynamicgo/mesh/libp2p/repo"
	"github.com/dynamicgo/slf4go"
	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	multiaddr "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

// Builtin protocols
var (
	ProtocolAdmin protocol.ID = "/mesh/admin/2.0.0"
)

// Node  .
type Node interface {
	Host() host.Host
	Detect(detector ProtocolDetector)
	Context() context.Context
	mesh.Network
}

type meshNode struct {
	slf4go.Logger            // logger
	host           host.Host // p2p host
	repo           repo.Repository
	ctx            context.Context
	cancel         context.CancelFunc
	bootstrapPeers []peerstore.PeerInfo
	config         config.Config
}

// New .
func New(ctx context.Context, config config.Config) (Node, error) {

	attachLibp2pLog()

	c, cancel := context.WithCancel(ctx)

	node := &meshNode{
		Logger: slf4go.Get("libp2p-mesh"),
		ctx:    c,
		cancel: cancel,
		config: config,
	}

	var err error
	node.repo, err = repo.Open(config.Get("mesh", "libp2p", "repo").String(".mesh"))

	if err != nil {
		node.ErrorF("open libp2p repo err: %s", err)
		return nil, err
	}

	if err := node.createHost(config); err != nil {
		node.ErrorF("open libp2p host err: %s", err)
		return nil, err
	}

	if err := node.loadBootstrapPeers(config); err != nil {
		node.ErrorF("open libp2p host err: %s", err)
		return nil, err
	}

	node.DebugF("node id: %s", node.host.ID().Pretty())

	return node, nil
}

func (node *meshNode) Context() context.Context {
	return node.ctx
}

func (node *meshNode) ID() string {
	return node.host.ID().Pretty()
}

func (node *meshNode) Stop() error {
	node.cancel()
	return nil
}

func (node *meshNode) createHost(config config.Config) error {

	var options []libp2p.Option

	var laddrs []string

	err := config.Get("mesh", "libp2p", "laddr").Scan(&laddrs)

	if err != nil {
		return err
	}

	if len(laddrs) == 0 {
		laddrs = []string{
			"/ip4/0.0.0.0/tcp/9000",
		}
	}

	options = append(options, libp2p.ListenAddrStrings(laddrs...))

	if !config.Get("mesh", "libp2p", "nat").Bool(false) {
		options = append(options, libp2p.NATPortMap())
	}

	privateKey, err := node.repo.PrivateKey()

	if err != nil {
		return err
	}

	options = append(options, libp2p.Identity(privateKey))

	privateKey = nil

	host, err := libp2p.New(node.ctx, options...)

	if err != nil {
		return err
	}

	node.host = host

	return nil
}

func (node *meshNode) loadBootstrapPeers(config config.Config) (err error) {

	var addrs []string

	err = config.Get("mesh", "libp2p", "bootstrap", "peers").Scan(&addrs)

	if err != nil {
		return
	}

	for _, a := range addrs {
		peerAddr, err := addr.Parse(a)

		if err != nil {
			node.WarnF("invalid peer addr :%s", a)
			continue
		}

		id, err := peer.IDB58Decode(peerAddr.ID)

		if err != nil {
			return err
		}

		node.bootstrapPeers = append(node.bootstrapPeers, peerstore.PeerInfo{
			ID:    id,
			Addrs: []multiaddr.Multiaddr{peerAddr.Addr},
		})
	}

	return
}

func (node *meshNode) Addrs() (addrs []string) {
	for _, a := range node.host.Addrs() {

		if manet.IsIPLoopback(a) {
			continue
		}

		peerAddr := &addr.PeerAddr{
			ID:   node.host.ID().Pretty(),
			Addr: a,
		}

		addrs = append(addrs, peerAddr.String())
	}
	return
}

func (node *meshNode) Host() host.Host {
	return node.host
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func randomSubsetOfPeers(in []peerstore.PeerInfo, max int) []peerstore.PeerInfo {
	n := min(max, len(in))
	var out []peerstore.PeerInfo
	for _, val := range rand.Perm(len(in)) {
		out = append(out, in[val])
		if len(out) >= n {
			break
		}
	}
	return out
}

func (node *meshNode) name() string {
	return node.host.ID().Pretty()
}

func (node *meshNode) boostrapOnce(config config.Config) {

	if len(node.bootstrapPeers) == 0 {
		node.DebugF("[%s] skip boostrap round with zero bootstrap peers", node.name())
		return
	}

	node.DebugF("[%s] start a new boostrap round ...", node.name())

	connected := node.host.Network().Peers()

	dails := config.Get("mesh", "libp2p", "minipeers").Int(10) - len(connected)

	if dails <= 0 {
		node.DebugF("[%s] living node %d, boostrap skippped", node.name(), len(connected))
		return
	}

	node.DebugF("[%s] connected peers %d ,try connect more peers", node.name(), len(connected))

	var notConnected []peerstore.PeerInfo
	peers := node.bootstrapPeers

	if len(peers) == 0 {
		node.WarnF("[%s] bootstrap peers table is zero", node.name())
	}

	for _, p := range peers {
		if node.host.Network().Connectedness(p.ID) != inet.Connected {
			notConnected = append(notConnected, p)
		}
	}

	notConnected = randomSubsetOfPeers(notConnected, dails)

	var wg sync.WaitGroup

	for _, peerInfo := range notConnected {
		wg.Add(1)
		go node.booststrapConnect(config, &wg, peerInfo)
	}

	wg.Wait()

	node.DebugF("[%s] end a boostrap round", node.name())
}

func (node *meshNode) booststrapConnect(config config.Config, wg *sync.WaitGroup, peerInfo peerstore.PeerInfo) {

	ctx, cancel := context.WithTimeout(node.ctx, config.Get("mesh", "libp2p", "timeout").Duration(time.Second*10))

	defer cancel()

	defer wg.Done()

	node.host.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)

	if err := node.host.Connect(ctx, peerInfo); err != nil {
		node.ErrorF("[%s] connect to %s err %s", node.name(), peerInfo.ID.String(), err)
	}
}

func (node *meshNode) Listen(serviceName string) (net.Listener, error) {
	return listener.Listen(node.ctx, node.host, protocol.ID(serviceName)), nil
}

func (node *meshNode) Dial(remote *mesh.Peer, serviceName string, timeout time.Duration) (net.Conn, error) {

	host := node.host

	id, err := peer.IDB58Decode(remote.ID)

	if err != nil {
		return nil, err
	}

	peerInfo := peerstore.PeerInfo{
		ID:    id,
		Addrs: remote.Addrs,
	}

	host.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)

	ctx, cancel := context.WithTimeout(node.ctx, timeout)

	defer cancel()

	err = host.Connect(ctx, peerInfo)

	if err != nil {
		node.ErrorF("[%s] ensure connect to remote peer %s err %s", node.name(), remote.ID, err)
		return nil, err
	}

	stream, err := host.NewStream(ctx, peerInfo.ID, protocol.ID(serviceName))

	if err != nil {
		return nil, err
	}

	conn := &netwrapper.StreamConn{Stream: stream}

	return conn, nil
}

func init() {
	mesh.RegisterNetwork("libp2p", func(ctx context.Context, config config.Config) (mesh.Network, error) {
		return New(ctx, config)
	})
}
