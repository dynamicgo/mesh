package mesh

import (
	"errors"

	"google.golang.org/grpc"
)

// Builtin errors
var (
	ErrServiceNotFound = errors.New("service not found")
	ErrNetworkClosed   = errors.New("underground network closed")
)

// Agent microsrevice mesh agent node
type Agent interface {
	Network() Network
	Stop() error
	RegisterService(name string, options ...grpc.ServerOption) (Service, error)
	FindService(name string, options ...grpc.DialOption) (*grpc.ClientConn, error)
}
