package sender

import (
	"code.google.com/p/gogoprotobuf/proto"
	"fmt"
	demitter "github.com/cloudfoundry/dropsonde/emitter"
	"github.com/cloudfoundry/dropsonde/events"
	"github.com/cloudfoundry/dropsonde/signature"
	"log"
)

type signedEventEmitter struct {
	innerEmitter demitter.ByteEmitter
	origin       string
	sharedSecret string
}

func NewSignedEventEmitter(targetAddress, origin, sharedSecret string) (e *signedEventEmitter, err error) {
	byteEmitter, err := demitter.NewUdpEmitter(targetAddress)
	if err != nil {
		return
	}
	e = &signedEventEmitter{byteEmitter, origin, sharedSecret}
	return
}

func (e *signedEventEmitter) Emit(event events.Event) error {
	envelope, err := demitter.Wrap(event, e.origin)
	if err != nil {
		return fmt.Errorf("Wrap: %v", err)
	}
	// for testing
	log.Printf("The envelope is %v", envelope)
	data, err := proto.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("Marshal: %v", err)
	}
	// for testing
	log.Printf("The data is %v", data)
	signedData := signature.SignMessage(data, []byte(e.sharedSecret))
	// for testing
	log.Printf("The signedData is %v", signedData)
	return e.innerEmitter.Emit(signedData)
}

func (e *signedEventEmitter) Close() {
	e.innerEmitter.Close()
}
