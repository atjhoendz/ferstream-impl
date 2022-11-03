package subscriber

import (
	"context"

	"github.com/kumparan/ferstream"
	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"
)

func getMsgFromFerstream(payload ferstream.MessageParser) (msg *ferstream.NatsEventMessage, err error) {
	logger := log.WithFields(log.Fields{
		"payload": utils.Dump(payload),
	})

	msg, ok := payload.(*ferstream.NatsEventMessage)
	if !ok {
		err = ferstream.ErrCastingPayloadToStruct
		logger.Error(err)
		return
	}

	return msg, nil
}

func createJetStreamEventHandler(streamName string, fn func(ctx context.Context, msg *ferstream.NatsEventMessage) error) ferstream.MessageHandler {
	return func(payload ferstream.MessageParser) (err error) {
		logger := log.WithFields(log.Fields{
			"payload": utils.Dump(payload),
		})

		msg, err := getMsgFromFerstream(payload)
		if err != nil {
			logger.Error(err)
			return err
		}

		return fn(context.Background(), msg)
	}
}
