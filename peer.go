package mesh

import (
	"errors"
	"fmt"
	"strings"

	multiaddr "github.com/multiformats/go-multiaddr"
)

// PINWE .
const (
	PINWE = 0x11BC
)

// Errors
var (
	ErrFormat = errors.New("invalid peer address format")
)

func init() {

	multiaddr.AddProtocol(multiaddr.Protocol{
		Code:       PINWE,
		Size:       multiaddr.LengthPrefixedVarSize,
		Name:       "mesh",
		VCode:      multiaddr.CodeToVarint(PINWE),
		Path:       false,
		Transcoder: multiaddr.TranscoderP2P,
	})
}

// Peer peer info struct
type Peer struct {
	ID    string
	Addrs []multiaddr.Multiaddr
}

// MeshAddrs convert peer to mesh addresses
func (peer *Peer) MeshAddrs() (addrs []string) {

	for _, addr := range peer.Addrs {
		suffix, err := multiaddr.NewMultiaddr(fmt.Sprintf("/mesh/%s", peer.ID))

		if err != nil {
			println("========", err.Error())
		}

		addrs = append(addrs, addr.Encapsulate(suffix).String())
	}

	return
}

// ParseMeshAddr parse mesh address to peer object
func ParseMeshAddr(name string) (*Peer, error) {
	m, err := multiaddr.NewMultiaddr(name)

	if err != nil {
		return nil, err
	}

	msplit := multiaddr.Split(m)
	if len(msplit) < 1 {
		return nil, ErrFormat
	}

	ipfspart := msplit[len(msplit)-1] // last part
	if ipfspart.Protocols()[0].Code != PINWE {
		return nil, ErrFormat
	}

	peerIDParts := strings.Split(ipfspart.String(), "/")
	peerIDStr := peerIDParts[len(peerIDParts)-1]

	return &Peer{
		ID:    peerIDStr,
		Addrs: []multiaddr.Multiaddr{multiaddr.Join(msplit[:len(msplit)-1]...)},
	}, nil
}

// MergePeers merge peers with same id
func MergePeers(peers []*Peer) []*Peer {
	peersMap := make(map[string]*Peer)

	for _, p := range peers {
		peer, ok := peersMap[p.ID]

		if !ok {
			peer = p
			peersMap[p.ID] = p
		}

		peer.Addrs = append(peer.Addrs, p.Addrs...)
	}

	peers = []*Peer{}

	for _, peer := range peersMap {
		peers = append(peers, peer)
	}

	return peers
}

// AddrsToPeers .
func AddrsToPeers(addrs []string) (peers []*Peer, err error) {

	for _, addr := range addrs {
		peer, err := ParseMeshAddr(addr)
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}

	return MergePeers(peers), nil
}
