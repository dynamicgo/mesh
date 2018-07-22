package hub

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/mesh/proto"
	"github.com/dynamicgo/slf4go"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
)

type serviceHub struct {
	slf4go.Logger
	client *redis.Client
}

func new(config config.Config) (proto.ServiceHubServer, error) {
	hub := &serviceHub{
		Logger: slf4go.Get("servicehub"),
	}

	hub.client = redis.NewClient(&redis.Options{
		Addr:     config.Get("mesh", "hub", "redis", "addrs").String(":6379"),
		Password: config.Get("mesh", "hub", "redis", "password").String(""),
	})

	return hub, nil
}

func (hub *serviceHub) Lookup(ctx context.Context, request *proto.LockupRequest) (*proto.LookupResponse, error) {

	hub.DebugF("lookup service[%s]", request.Name)

	addrs, err := hub.client.SMembers(request.Name).Result()

	if err != nil {
		return nil, err
	}

	return &proto.LookupResponse{Addrs: addrs}, nil
}

func (hub *serviceHub) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	hub.DebugF("register service[%s] from peer with addrs[%s]", request.Name, strings.Join(request.Addrs, ","))

	var members []interface{}

	for _, addr := range request.Addrs {
		members = append(members, addr)
	}

	err := hub.client.SAdd(request.Name, members...).Err()

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "call redis err: %s", err)
	}

	return &proto.RegisterResponse{}, nil
}
