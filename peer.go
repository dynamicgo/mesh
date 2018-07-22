package mesh

import (
	"github.com/dynamicgo/mesh/addr"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// Peer peer info struct
type Peer struct {
	ID    string
	Addrs []multiaddr.Multiaddr
}

// AddrsToPeers .
func AddrsToPeers(addrs []string) (peers []*Peer, err error) {

	peersMap := make(map[string]*Peer)

	for _, a := range addrs {
		peerAddr, err := addr.Parse(a)

		if err != nil {
			return nil, err
		}

		peer, ok := peersMap[peerAddr.ID]

		if !ok {
			peer = &Peer{
				ID: peerAddr.ID,
			}
			peersMap[peerAddr.ID] = peer
		}

		peer.Addrs = append(peer.Addrs, peerAddr.Addr)
	}

	for _, peer := range peersMap {
		peers = append(peers, peer)
	}

	return
}
