package sender

import (
	"github.com/amzfans/loggregator_agent/listener"
	"github.com/cloudfoundry/dropsonde"
	demitter "github.com/cloudfoundry/dropsonde/emitter"
	"github.com/cloudfoundry/dropsonde/logs"
	"strconv"
)

const (
	origin_prefix = "app_loggregator_agent"
	source_name   = "APP"
)

type LogSender struct {
	appId         string
	index         uint64
	targetAddress string
	sharedSecret  string
	logListener   *listener.LogListener
	emitter       demitter.EventEmitter
	StopChan      chan bool
}

func NewLogSender(appId, targetAddress, sharedSecret string, index uint64,
	logListener *listener.LogListener) (lsender *LogSender, err error) {
	lsender = &LogSender{appId, index, targetAddress, sharedSecret, logListener, nil, make(chan bool, 1)}
	lsender.emitter, err = NewSignedEventEmitter(targetAddress, origin_prefix+"/"+appId+"/"+strconv.FormatUint(index, 10), sharedSecret)
	// initialize the dropsonde with emitter
	dropsonde.InitializeWithEmitter(lsender.emitter)
	return
}

func (lsender *LogSender) Start() {
	go func() {
		defer lsender.Stop()
		logs.ScanLogStream(lsender.appId, source_name,
			strconv.FormatUint(lsender.index, 10), lsender.logListener.StdoutConn)
	}()

	go func() {
		defer lsender.Stop()
		logs.ScanErrorLogStream(lsender.appId, source_name,
			strconv.FormatUint(lsender.index, 10), lsender.logListener.StderrConn)
	}()
}

func (lsender *LogSender) Stop() {
	lsender.StopChan <- true
	lsender.emitter.Close()
}
