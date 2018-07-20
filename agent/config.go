package agent

import (
	"context"
	"time"

	"github.com/dynamicgo/slf4go"

	"github.com/dynamicgo/go-config/source"
	proto "github.com/dynamicgo/go-config/source/grpc/proto"
)

type grpcSource struct {
	slf4go.Logger
	service string
	client  proto.SourceClient
}

func (agent *agentImpl) newGrpcSource(service string, client proto.SourceClient) source.Source {
	return &grpcSource{
		Logger:  agent.Logger,
		service: service,
		client:  client,
	}
}

func (g *grpcSource) Read() (*source.ChangeSet, error) {

	rsp, err := g.client.Read(context.Background(), &proto.ReadRequest{
		Path: g.service,
	})
	if err != nil {
		return nil, err
	}
	return toChangeSet(rsp.ChangeSet), nil
}

func (g *grpcSource) Watch() (source.Watcher, error) {

	rsp, err := g.client.Watch(context.Background(), &proto.WatchRequest{
		Path: g.service,
	})
	if err != nil {
		return nil, err
	}
	return newWatcher(rsp)
}

func (g *grpcSource) String() string {
	return "grpc"
}

func toChangeSet(c *proto.ChangeSet) *source.ChangeSet {
	return &source.ChangeSet{
		Data:      c.Data,
		Checksum:  c.Checksum,
		Format:    c.Format,
		Timestamp: time.Unix(c.Timestamp, 0),
		Source:    c.Source,
	}
}

type watcher struct {
	stream proto.Source_WatchClient
}

func newWatcher(stream proto.Source_WatchClient) (*watcher, error) {
	return &watcher{
		stream: stream,
	}, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	rsp, err := w.stream.Recv()
	if err != nil {
		return nil, err
	}
	return toChangeSet(rsp.ChangeSet), nil
}

func (w *watcher) Stop() error {
	return w.stream.CloseSend()
}
