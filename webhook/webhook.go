package webhook

import (
	"os"
	"os/signal"
	"sync"

	"net/http"

	config "github.com/chrusty/prometheus_webhook_snmptrapper/config"
	types "github.com/chrusty/prometheus_webhook_snmptrapper/types"

	logrus "github.com/sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{"logger": "Webhook-server"})
	myConfig config.Config
)

func init() {
	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)
}

func Run(myConfigFromMain config.Config, alertsChannel chan types.Alert, waitGroup *sync.WaitGroup) {

	log.WithFields(logrus.Fields{"address": myConfigFromMain.WebhookAddress}).Info("Starting the Webhook server")

	// Populate the config:
	myConfig = myConfigFromMain

	// Set up a channel to handle shutdown:
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)

	// Listen for webhooks:
	http.ListenAndServe(myConfig.WebhookAddress, &WebhookHandler{AlertsChannel: alertsChannel})

	// Wait for shutdown:
	for {
		select {
		case <-signals:
			log.Info("Shutting down the Webhook server")

			// Tell main() that we're done:
			waitGroup.Done()
			return
		}
	}

}
