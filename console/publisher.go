package console

import (
	"os"
	"os/signal"
	"time"

	"github.com/atjhoendz/ferstream-impl/config"
	"github.com/atjhoendz/ferstream-impl/usecase"
	"github.com/jaswdr/faker"
	"github.com/kumparan/ferstream"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var publisherCmd = &cobra.Command{
	Use:   "publisher",
	Short: "start publisher",
	Run:   publisher,
}

func init() {
	publisherCmd.Flags().String("type", "once", "publish type")

	rootCmd.AddCommand(publisherCmd)
}

func publisher(cmd *cobra.Command, _ []string) {
	publishType, err := cmd.Flags().GetString("type")
	if err != nil {
		log.Fatal(err)
	}

	helloUsecase := usecase.NewHelloUsecase()

	jsClients := []ferstream.JetStreamRegistrar{helloUsecase}
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

	switch publishType {
	case "once":
		err := publishOneMessage(helloUsecase)
		if err != nil {
			log.Error(err)
			return
		}
	case "multiple":
		ticker := time.NewTicker(2 * time.Second)
		go publishMultipleMessage(helloUsecase, done, ticker)
		defer ticker.Stop()
	default:
		log.Errorf("argument %s invalid", publishType)
	}

	<-done
}

func publishOneMessage(helloUsecase *usecase.HelloUsecase) error {
	name := faker.New().Person().Name()
	err := helloUsecase.SayHello(name)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("sent hello message to %s", name)

	return nil
}

func publishMultipleMessage(helloUsecase *usecase.HelloUsecase, done <-chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			name := faker.New().Person().Name()
			err := helloUsecase.SayHello(name)
			if err != nil {
				log.Error(err)
			}

			log.Infof("sent hello message to %s", name)
		}
	}
}
