package mesh

import (
	config "github.com/dynamicgo/go-config"
	"google.golang.org/grpc"
)

// Builtin protocols
var (
	ServiceHub    = "/dynamicgo/mesh/hub/1.0.0"
	ConfigService = "/dynamicgo/mesh/config/1.0.0"
)

// ServiceMain service main entry
type ServiceMain func(agent Agent, server *grpc.Server, config config.Config) error

// Service .
type Service interface {
	Run(main ServiceMain, options ...ServiceOption) error
}
