package mesh

import (
	"github.com/dynamicgo/go-config/source"
)

// FindOption .
type FindOption func()

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
