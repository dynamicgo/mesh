package mesh

import (
	"github.com/dynamicgo/slf4go"
	"google.golang.org/grpc/grpclog"
)

type grpcLogger struct {
	logger slf4go.Logger
}

func (log *grpcLogger) Info(args ...interface{}) {
	log.logger.Info(args...)
}
func (log *grpcLogger) Infoln(args ...interface{}) {
	log.logger.Info(args...)
}
func (log *grpcLogger) Infof(format string, args ...interface{}) {
	log.logger.InfoF(format, args...)
}
func (log *grpcLogger) Warning(args ...interface{}) {
	log.logger.Warn(args...)
}
func (log *grpcLogger) Warningln(args ...interface{}) {
	log.logger.Warn(args...)
}
func (log *grpcLogger) Warningf(format string, args ...interface{}) {
	log.logger.WarnF(format, args...)
}
func (log *grpcLogger) Error(args ...interface{}) {
	log.logger.Error(args...)
}
func (log *grpcLogger) Errorln(args ...interface{}) {
	log.logger.Error(args...)
}
func (log *grpcLogger) Errorf(format string, args ...interface{}) {
	log.logger.ErrorF(format, args...)
}
func (log *grpcLogger) Fatal(args ...interface{}) {
	log.logger.Fatal(args...)
}
func (log *grpcLogger) Fatalln(args ...interface{}) {
	log.logger.Warn(args...)
}
func (log *grpcLogger) Fatalf(format string, args ...interface{}) {
	log.logger.FatalF(format, args...)
}
func (log *grpcLogger) V(l int) bool {
	return true
}

func init() {
	logger := &grpcLogger{logger: slf4go.Get("grpc")}
	logger.logger.SourceCodeLevel(5)
	grpclog.SetLoggerV2(logger)
}
