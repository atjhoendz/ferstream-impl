package usecase

import (
	"fmt"
	"time"

	"github.com/atjhoendz/ferstream-impl/event"
	"github.com/kumparan/ferstream"
	"github.com/kumparan/go-utils"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type Hello struct {
	js ferstream.JetStream
}

type HelloUsecase interface {
	RegisterNATSJetStream(js ferstream.JetStream)
	SayHello(msg string) error
}

func NewHelloUsecase() HelloUsecase {
	return &Hello{}
}

func (h *Hello) RegisterNATSJetStream(js ferstream.JetStream) {
	h.js = js
}

func (h *Hello) InitStream() error {
	_, err := h.js.AddStream(&nats.StreamConfig{
		Name:     event.HelloStreamName,
		Subjects: []string{event.HelloStreamSubjectAll},
		MaxAge:   24 * time.Hour,
		Storage:  nats.FileStorage,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (h *Hello) SayHello(msg string) error {
	eventMsg := ferstream.NewNatsEventMessage().
		WithEvent(
			&ferstream.NatsEvent{
				ID:      utils.GenerateID(),
				UserID:  111,
				Subject: event.HelloStreamSubjectCreate,
			}).
		WithBody(fmt.Sprintf("hello %s", msg))

	msgByte, err := eventMsg.Build()
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = h.js.Publish(event.HelloStreamSubjectCreate, msgByte)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
