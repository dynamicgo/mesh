package hub

import (
	"context"

	"github.com/dynamicgo/mesh/proto"
	"github.com/dynamicgo/slf4go"
)

type serviceHub struct {
	slf4go.Logger
}

func (hub *serviceHub) ServiceLookup(context.Context, *proto.LockupRequest) (*proto.LookupResponse, error) {
	return nil, nil
}

func (hub *serviceHub) Register(context.Context, *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return nil, nil
}
