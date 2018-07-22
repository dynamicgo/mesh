package agent

import (
	"fmt"
	"testing"

	"github.com/dynamicgo/orm"

	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/go-config/source/memory"
	"github.com/dynamicgo/mesh"
	_ "github.com/dynamicgo/mesh-libp2p-network"
	"github.com/dynamicgo/mesh/service/configservice"
	"github.com/dynamicgo/mesh/service/hub"
	_ "github.com/mattn/go-sqlite3"
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
		"libp2p":{
			"laddr":[
				"/ip4/0.0.0.0/tcp/8000"
			],
			"repo":".test/agent"
		}
	}
}
`

var cdata = `
{
	"mesh":{
		"hub":{
			"peers":[
				"/ip4/127.0.0.1/tcp/7000/mesh/%s"
			]
		},
		"libp2p":{
			"laddr":[
				"/ip4/0.0.0.0/tcp/9000"
			],
			"repo":".test/cs"
		}
	}
}
`

var node mesh.Agent
var servicehub mesh.Agent
var configService mesh.Agent

func init() {
	config := config.NewConfig(config.WithSource(memory.NewSource(memory.WithData(hubconfigdata))))

	driver := config.Get("mesh", "configservice", "database", "driver").String("sqlite3")
	source := config.Get("mesh", "configservice", "database", "source").String(".test/configservice.db")

	db, err := xorm.NewEngine(driver, source)

	if err != nil {
		panic(err)
	}

	if err := orm.Sync(db); err != nil {
		panic(err)
	}

	servicehub, err = New(config)

	if err != nil {
		panic(err.Error())
	}

	service, err := servicehub.RegisterService(mesh.ProtocolServiceHub)

	if err != nil {
		panic(err)
	}

	go func() {
		err := service.Run(hub.Main, mesh.NoRemoteConfig())
		if err != nil {
			panic(err)
		}
	}()

	// go service.Run(hub.Main, mesh.NoRemoteConfig())
}

func init() {
	configdata := []byte(fmt.Sprintf(cdata, servicehub.Network().ID()))

	config := config.NewConfig(config.WithSource(memory.NewSource(memory.WithData(configdata))))

	var err error
	configService, err = New(config)

	if err != nil {
		panic(err.Error())
	}

	service, err := configService.RegisterService(mesh.ProtocolConfigService)

	if err != nil {
		panic(err)
	}

	go func() {
		err := service.Run(configservice.Main, mesh.NoRemoteConfig())

		if err != nil {
			// println(err)
			panic(err)
		}
	}()
}

func init() {

	configdata := []byte(fmt.Sprintf(clientconfigdata, servicehub.Network().ID()))

	// configdata := []byte(fmt.Sprintf(clientconfigdata, "", ""))

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

	err = service.Run(func(mesh.Agent, *grpc.Server, config.Config) error {
		return nil
	})

	require.NoError(t, err)
}
