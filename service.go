package mesh

import (
	config "github.com/dynamicgo/go-config"
	"github.com/dynamicgo/go-config/source"
	"google.golang.org/grpc"
)

// Builtin protocols
var (
	ServiceHub    = "/dynamicgo/mesh/hub/1.0.0"
	ConfigService = "/dynamicgo/mesh/config/1.0.0"
)

// ServiceConfig .
type ServiceConfig struct {
	ConfigSources []source.Source
	RemoteConfig  bool
}

// ServiceOption .
type ServiceOption func(config *ServiceConfig)

// NoRemoteConfig ignore load config from admin
func NoRemoteConfig() ServiceOption {
	return func(config *ServiceConfig) {
		config.RemoteConfig = false
	}
}

// WithConfig load service with input config ,multi WithConfig option will auto merged
func WithConfig(source source.Source) ServiceOption {
	return func(serviceConfig *ServiceConfig) {
		serviceConfig.ConfigSources = append(serviceConfig.ConfigSources, source)
	}
}

// ServiceMain service main entry
type ServiceMain func(agent Agent, server *grpc.Server, config config.Config) error

// Service .
type Service interface {
	Run(main ServiceMain, options ...ServiceOption) error
}
