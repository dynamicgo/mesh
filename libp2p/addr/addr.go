package addr

import (
	"errors"
	"fmt"
	"strings"

	peer "github.com/libp2p/go-libp2p-peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// PINWE .
const (
	PINWE = 0x11BC
)

func init() {

	multiaddr.AddProtocol(multiaddr.Protocol{
		Code:       PINWE,
		Size:       multiaddr.LengthPrefixedVarSize,
		Name:       "mesh",
		VCode:      multiaddr.CodeToVarint(PINWE),
		Path:       false,
		Transcoder: multiaddr.TranscoderIPFS,
	})
}

// Errors
var (
	ErrFormat = errors.New("invalid address format")
)

// PeerAddr inwe peer address
type PeerAddr struct {
	ID   peer.ID
	Addr multiaddr.Multiaddr
}

func (addr *PeerAddr) String() string {
	suffix, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/mesh/%s", addr.ID.Pretty()))

	return addr.Addr.Encapsulate(suffix).String()
}

// Parse parse string to peer address oject
func Parse(name string) (*PeerAddr, error) {

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

	// make sure 'ipfs id' parses as a peer.ID
	peerIDParts := strings.Split(ipfspart.String(), "/")
	peerIDStr := peerIDParts[len(peerIDParts)-1]
	id, err := peer.IDB58Decode(peerIDStr)

	if err != nil {
		return nil, err
	}

	return &PeerAddr{
		ID:   id,
		Addr: multiaddr.Join(msplit[:len(msplit)-1]...),
	}, nil
}
