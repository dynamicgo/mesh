package libp2p

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/go-config/source/memory"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
)

var ctrlconfigdata = []byte(`
{
	"mesh":{
		"libp2p":{
			"laddr":[
				"/ip4/0.0.0.0/tcp/7000"
			],
			"repo":".test/ctrl"
		}
	}
}
`)

var clientconfigdata = `
{
	"mesh":{
		"libp2p":{
			"laddr":[
				"/ip4/0.0.0.0/tcp/8000"
			],
			"repo":".test/agent",
			"bootstrap":{
				"peers":[
					"/ip4/127.0.0.1/tcp/7000/mesh/%s"
				]
			}
		}
	}
}
`

var ctrlnode Node

func registerCtr(node Node) {
	node.Host().SetStreamHandler(ProtocolAdmin, func(stream inet.Stream) {

	})
}

func init() {
	ctrlconfig := config.NewConfig(config.WithSource(memory.NewSource(memory.WithData(ctrlconfigdata))))

	var err error

	ctrlnode, err = New(context.Background(), ctrlconfig)

	if err != nil {
		panic(err)
	}

	registerCtr(ctrlnode)

}

type detector struct {
	updated chan peer.ID
}

func newDetector() *detector {
	return &detector{
		updated: make(chan peer.ID),
	}
}

func (d *detector) Protocols() []protocol.ID {
	return []protocol.ID{
		ProtocolAdmin,
	}
}

func (d *detector) Update(protocol []protocol.ID, id peer.ID) {
	d.updated <- id
}

func (d *detector) Removed(id peer.ID) {

}

func (d *detector) Name() string {
	return "agent-detector"
}

func TestDetector(t *testing.T) {
	configdata := []byte(fmt.Sprintf(clientconfigdata, ctrlnode.Host().ID().Pretty()))

	config := config.NewConfig(config.WithSource(memory.NewSource(memory.WithData(configdata))))

	agent, err := New(context.Background(), config)

	require.NoError(t, err)

	de := newDetector()

	agent.Detect(de)

	<-de.updated

	agent.Stop()
}
