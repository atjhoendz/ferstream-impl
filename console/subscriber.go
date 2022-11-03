package console

import (
	"os"
	"os/signal"

	"github.com/atjhoendz/ferstream-impl/config"
	"github.com/atjhoendz/ferstream-impl/subscriber"
	"github.com/kumparan/ferstream"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var subscriberCmd = &cobra.Command{
	Use:   "subscriber",
	Short: "start subscriber",
	Run:   initSubscriber,
}

func init() {
	subscriberCmd.Flags().String("type", "queue-subscribe", "subscription type")

	rootCmd.AddCommand(subscriberCmd)
}

func initSubscriber(cmd *cobra.Command, _ []string) {
	subscriptionType, err := cmd.Flags().GetString("type")
	if err != nil {
		log.Fatal(err)
	}

	var jsClients []ferstream.JetStreamRegistrar

	switch subscriptionType {
	case "queue-subscribe":
		queueSub := subscriber.NewQueueSub()
		jsClients = append(jsClients, queueSub)
	case "subscribe":
		// TODO
	default:
		log.Error("argument invalid")
		return
	}

	js, err := ferstream.NewNATSConnection(
		config.DefaultNATSJSHost,
		jsClients,
		nats.RetryOnFailedConnect(config.DefaultNATSJSRetryOnFailedConnect),
		nats.MaxReconnects(config.DefaultNATSJSMaxReconnect),
		nats.ReconnectWait(config.DefaultNATSJSReconnectWait),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	signalChan := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)

	js.GetNATSConnection().SetClosedHandler(func(c *nats.Conn) {
		done <- true
	})

	go func(ferstream.JetStream) {
		for range signalChan {
			log.Println("received an interrupt...")
			log.Println("draining nats connection...")

			// Drain the connection, which will close it when done.
			ferstream.SafeClose(js)
			log.Println("exiting...")
			return
		}
	}(js)

	<-done
}
