package hub

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"

	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/mesh/proto"
	"github.com/dynamicgo/slf4go"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
)

type serviceHub struct {
	slf4go.Logger
	client  *redis.Client
	timeout time.Duration
}

func new(config config.Config) (proto.ServiceHubServer, error) {
	hub := &serviceHub{
		Logger:  slf4go.Get("servicehub"),
		timeout: config.Get("mesh", "hub", "timeout").Duration(time.Second * 60),
	}

	hub.client = redis.NewClient(&redis.Options{
		Addr:     config.Get("mesh", "hub", "redis", "addrs").String(":6379"),
		Password: config.Get("mesh", "hub", "redis", "password").String(""),
	})

	return hub, nil
}

func (hub *serviceHub) Lookup(ctx context.Context, request *proto.LockupRequest) (*proto.LookupResponse, error) {

	hub.DebugF("lookup service[%s]", request.Name)

	keys, err := hub.client.SMembers(request.Name).Result()

	if err != nil {
		return nil, err
	}

	hub.DebugF("lookup service[%s] get keys [%s]", request.Name, strings.Join(keys, ","))

	var addrs []string

	for _, key := range keys {
		val, err := hub.client.Get(key).Result()

		if err != nil {
			if err != redis.Nil {
				return nil, fmt.Errorf("get key[%s] from redis err: %s", key, err)
			}

			if err := hub.client.SRem(request.Name, key).Err(); err != nil {
				return nil, fmt.Errorf("remove key[%s] from redis err: %s", key, err)
			}
		}

		addrs = append(addrs, strings.Split(val, ",")...)
	}

	return &proto.LookupResponse{Addrs: addrs}, nil
}

func (hub *serviceHub) Register(ctx context.Context, request *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	peer, ok := peer.FromContext(ctx)

	if !ok {
		return nil, fmt.Errorf("can't find peer")
	}

	hub.DebugF("register service[%s] from peer with addrs[%s]", request.Name, strings.Join(request.Addrs, ","))

	addrs := strings.Join(request.Addrs, ",")

	key := fmt.Sprintf("%s||%s", request.Name, peer.Addr.String())

	hub.DebugF("generate key %s", key)

	err := hub.client.Set(key, addrs, hub.timeout).Err()

	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "call redis err: %s", err)
	}

	hub.client.SAdd(request.Name, key)

	return &proto.RegisterResponse{}, nil
}
