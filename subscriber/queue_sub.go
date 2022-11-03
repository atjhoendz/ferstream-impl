package subscriber

import (
	"context"

	"github.com/atjhoendz/ferstream-impl/config"
	"github.com/atjhoendz/ferstream-impl/event"
	"github.com/kumparan/ferstream"
	"github.com/kumparan/go-utils"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

type QueueSub struct {
	js ferstream.JetStream
}

const queueGroup = "hello-queue-group"

func NewQueueSub() *QueueSub {
	return &QueueSub{}
}

func (q *QueueSub) RegisterNATSJetStream(js ferstream.JetStream) {
	q.js = js
}

func (q *QueueSub) SubscribeJetStreamEvent() error {
	eventHandler := createJetStreamEventHandler(
		event.HelloStreamName,
		func(ctx context.Context, msg *ferstream.NatsEventMessage) error {
			log.Infof("msg: %s", msg.Body)
			return nil
		})

	err := utils.Retry(config.DefaultNATSJSSubscribeRetryAttempts, config.DefaultNATSJSSubscribeRetryInterval, func() error {
		_, err := q.js.QueueSubscribe(
			event.HelloStreamSubjectAll,
			queueGroup,
			ferstream.NewNATSMessageHandler(
				new(ferstream.NatsEventMessage),
				config.DefaultNATSJSRetryAttempts,
				config.DefaultNATSJSRetryInterval,
				eventHandler,
				nil,
			),
			nats.ManualAck(),
			nats.Durable(config.NATSDurableID),
		)
		return err
	})
	if err != nil {
		log.Error(err)
	}

	return err
}
