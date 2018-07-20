package agent

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/go-config/source/memory"
	"github.com/dynamicgo/mesh"
	_ "github.com/dynamicgo/mesh/libp2p"
	"github.com/dynamicgo/mesh/service/configservice"
	"github.com/dynamicgo/mesh/service/hub"
)

var hubconfigdata = []byte(`
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
		"hub":{
			"peers":[
				"/ip4/127.0.0.1/tcp/7000/mesh/%s"
			]
		},
		"config":{
			"peers":[
				"/ip4/127.0.0.1/tcp/7000/mesh/%s"
			]
		},
		"libp2p":{
			"laddr":[
				"/ip4/0.0.0.0/tcp/8000"
			],
			"repo":".test/agent"
		}
	}
}
`

var node mesh.Agent
var servicehub mesh.Agent

func init() {
	config := config.NewConfig(config.WithSource(memory.NewSource(memory.WithData(hubconfigdata))))
	var err error
	servicehub, err = New(config)

	if err != nil {
		panic(err.Error())
	}

	service, err := servicehub.RegisterService(mesh.ServiceHub)

	if err != nil {
		panic(err)
	}

	go service.Run(hub.Main, mesh.NoRemoteConfig())

	service, err = servicehub.RegisterService(mesh.ConfigService)

	if err != nil {
		panic(err)
	}

	go service.Run(configservice.Main, mesh.NoRemoteConfig())
}

func init() {

	configdata := []byte(fmt.Sprintf(clientconfigdata, servicehub.Network().ID(), servicehub.Network().ID()))

	config := config.NewConfig(config.WithSource(memory.NewSource(memory.WithData(configdata))))

	var err error
	node, err = New(config)

	if err != nil {
		panic(err.Error())
	}
}

func TestRegisterService(t *testing.T) {
	service, err := node.RegisterService("test")

	require.NoError(t, err)

	err = service.Run(func(*grpc.Server, config.Config) error {
		return nil
	})

	require.NoError(t, err)
}
