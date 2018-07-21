package configservice

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/go-xorm/xorm"

	"github.com/dynamicgo/slf4go"

	config "github.com/dynamicgo/go-config"
	configgrpc "github.com/dynamicgo/go-config/source/grpc/proto"
)

type serviceImpl struct {
	slf4go.Logger
	orm *xorm.Engine
}

func new(config config.Config) (configgrpc.SourceServer, error) {

	driver := config.Get("mesh", "configservice", "database", "driver").String("sqlite3")
	source := config.Get("mesh", "configservice", "database", "source").String(".test/configservice.db")

	orm, err := xorm.NewEngine(driver, source)

	if err != nil {
		return nil, err
	}

	return &serviceImpl{
		Logger: slf4go.Get("configservice"),
		orm:    orm,
	}, nil
}
func (service *serviceImpl) Read(context context.Context, request *configgrpc.ReadRequest) (*configgrpc.ReadResponse, error) {

	var serviceConfig ServiceConfig

	_, err := service.orm.Where(`"path" = ?`, request.Path).Get(&serviceConfig)

	if err != nil {
		return nil, err
	}

	hash := md5.New()

	if _, err := hash.Write([]byte(serviceConfig.Content)); err != nil {
		return nil, err
	}

	response := &configgrpc.ReadResponse{
		ChangeSet: &configgrpc.ChangeSet{
			Data:      []byte(serviceConfig.Content),
			Checksum:  hex.EncodeToString(hash.Sum(nil)),
			Format:    "json",
			Source:    "configservice",
			Timestamp: time.Now().Unix(),
		},
	}

	return response, nil
}
func (service *serviceImpl) Watch(*configgrpc.WatchRequest, configgrpc.Source_WatchServer) error {
	return nil
}
