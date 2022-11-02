package usecase

import (
	"fmt"
	"time"

	"github.com/kumparan/ferstream"
	"github.com/kumparan/go-utils"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"gitlab.com/achun.armando/ferstream-impl/event"
)

type HelloUsecase struct {
	js ferstream.JetStream
}

func NewHelloUsecase() *HelloUsecase {
	return &HelloUsecase{}
}

func (h *HelloUsecase) RegisterNATSJetStream(js ferstream.JetStream) {
	h.js = js
}

func (h *HelloUsecase) InitStream() error {
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

func (h *HelloUsecase) SayHello(msg string) error {
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
