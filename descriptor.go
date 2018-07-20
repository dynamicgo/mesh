package mesh

// // ServiceDescriptor .
// type ServiceDescriptor struct {
// 	Name             string              // service unique name
// 	Balance          bool                // service balance
// 	GRPCServerOption []grpc.ServerOption // grpc server options
// }

// // ServiceOption .
// type ServiceOption func(descriptor *ServiceDescriptor)

// // WithName service name option
// func WithName(name string) ServiceOption {
// 	return func(config *ServiceDescriptor) {
// 		config.Name = name
// 	}
// }

// // WithGRPCServerOption service name option
// func WithGRPCServerOption(option grpc.ServerOption) ServiceOption {
// 	return func(config *ServiceDescriptor) {
// 		config.GRPCServerOption = append(config.GRPCServerOption, option)
// 	}
// }

// // MakeServiceDescriptor .
// func MakeServiceDescriptor(options ...ServiceOption) *ServiceDescriptor {
// 	serviceConfig := &ServiceDescriptor{}

// 	for _, option := range options {
// 		option(serviceConfig)
// 	}

// 	return serviceConfig
// }
