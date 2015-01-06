package listener

import (
	"github.com/cloudfoundry/dropsonde/events"
	"log"
	"net"
	"time"
)

const (
	stdout_socket_file = "stdout.sock"
	stderr_socket_file = "stderr.sock"
)

type LogListener struct {
	StdoutConn net.Conn
	StderrConn net.Conn
}

func NewLogListener() *LogListener {
	return &LogListener{}
}

func (ls *LogListener) Start() {
	go ls.createConnection(events.LogMessage_OUT)
	go ls.createConnection(events.LogMessage_ERR)
}

func (ls *LogListener) Stop() {
	if ls.StdoutConn != nil {
		ls.StdoutConn.Close()
	}

	if ls.StderrConn != nil {
		ls.StderrConn.Close()
	}
}

func connectToSocket(socketFile string) (conn net.Conn, err error) {
	for i := 0; i < 10; i++ {
		conn, err = net.Dial("unix", socketFile)
		if err == nil {
			log.Printf("Opened socket for %s.", socketFile)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return
}

func (ls *LogListener) createConnection(messageType events.LogMessage_MessageType) {
	var conn net.Conn
	var err error
	switch messageType {
	case events.LogMessage_OUT:
		conn, err = connectToSocket(stdout_socket_file)
		ls.StdoutConn = conn
	case events.LogMessage_ERR:
		conn, err = connectToSocket(stderr_socket_file)
		ls.StderrConn = conn
	default:
		log.Printf("Unknown messageType: %s.", messageType)
		return
	}

	if err != nil {
		log.Printf("Cannot open socket for %s as %s", messageType, err.Error())
		return
	}
}
