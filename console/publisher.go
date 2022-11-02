package console

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jaswdr/faker"
	"github.com/kumparan/ferstream"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/achun.armando/ferstream-impl/config"
	"gitlab.com/achun.armando/ferstream-impl/usecase"
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
	defer ferstream.SafeClose(js)

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	signalChan := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)

	go func(ferstream.JetStream) {
		for range signalChan {
			log.Println("received an interrupt...")
			log.Println("draining nats connection...")

			// Drain the connection, which will close it when done.
			ferstream.SafeClose(js)

			log.Println("exiting...")
			done <- true
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
		log.Error("argument invalid")
	}

	<-done
	fmt.Scanln()
}

// func initNATSJSClients(NATSJetStreamHost string, clients []ferstream.JetStreamRegistrar) (js ferstream.JetStream, err error) {
// 	natsConf := []nats.Option{
// 		nats.UseOldRequestStyle(),
// 		nats.RetryOnFailedConnect(config.DefaultNATSJSRetryOnFailedConnect),
// 		nats.MaxReconnects(config.DefaultNATSJSMaxReconnect),
// 		nats.ReconnectWait(config.DefaultNATSJSReconnectWait),
// 		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
// 			log.Errorf("NATS got error! Reason: %q\n", err)
// 		}),
// 		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
// 			log.Errorf("NATS got disconnected! Reason: %q\n", err)
// 		}),
// 		// nats.ReconnectHandler(ferstream.ReconnectHandler(clients)),
// 		nats.ReconnectHandler(func(c *nats.Conn) {
// 			log.Info("RECONNECT NICH")
// 		}),
// 		nats.ClosedHandler(func(nc *nats.Conn) {
// 			log.Errorf("NATS connection closed. Reason: %q\n", nc.LastError())
// 		}),
// 	}
// 	js, err = ferstream.NewNATSConnection(NATSJetStreamHost, natsConf...)
// 	if err != nil {
// 		log.Errorf("Failed to connect nats server. Reason: %q\n", err)
// 		return
// 	}

// 	err = ferstream.RegisterJetStreamClient(js, clients)
// 	if err != nil {
// 		log.Error(err)
// 	}

// 	return
// }

func publishOneMessage(helloUsecase *usecase.HelloUsecase) error {
	err := helloUsecase.SayHello("brader")
	if err != nil {
		log.Error(err)
		return err
	}

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
