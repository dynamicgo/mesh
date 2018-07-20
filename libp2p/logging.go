package libp2p

import (
	"strings"
	"sync"

	"github.com/dynamicgo/slf4go"
	logging "github.com/whyrusleeping/go-logging"
)

var libp2pLogOnce sync.Once

var libp2pLogFormat = logging.MustStringFormatter(
	"%{level:.4s}:%{message}",
)

type slf4goWriter struct {
	slf4go.Logger
}

func (writer slf4goWriter) Write(p []byte) (n int, err error) {

	messages := strings.SplitN(string(p), ":", 2)

	switch messages[0] {
	case "DEBU":
		writer.Debug(messages[1])
	case "INFO":
		writer.Info(messages[1])
	case "NOTI":
		writer.Info(messages[1])
	case "WARN":
		writer.Warn(messages[1])
	case "ERRO":
		writer.Error(messages[1])
	case "CRIT":
		writer.Error(messages[1])
	default:
		writer.Debug(messages[1])
	}

	return len(p), nil
}

func attachLibp2pLog() {

	libp2pLogOnce.Do(func() {
		logger := slf4go.Get("libp2p")

		logger.SourceCodeLevel(11)

		slf4goBackend := logging.NewLogBackend(&slf4goWriter{logger}, "", 0)

		logging.SetBackend(logging.NewBackendFormatter(slf4goBackend, libp2pLogFormat))
	})

}
