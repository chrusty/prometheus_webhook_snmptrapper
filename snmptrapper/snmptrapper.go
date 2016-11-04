package snmptrapper

import (
	"os"
	"os/signal"
	"sync"

	config "github.com/chrusty/prometheus_webhook_snmptrapper/config"
	types "github.com/chrusty/prometheus_webhook_snmptrapper/types"

	logrus "github.com/Sirupsen/logrus"
)

var (
	log      = logrus.WithFields(logrus.Fields{"logger": "SNMP-trapper"})
	myConfig config.Config
	trapOIDs types.TrapOIDs
)

func init() {
	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)

	// Configure which OIDs to use for the SNMP Traps:
	trapOIDs = types.TrapOIDs{
		TrapOID:      ".1.3.6.1.4.1.56.12.1.7",
		Component:    ".1.3.6.1.4.1.56.12.9.1.0",
		Message:      ".1.3.6.1.4.1.56.12.9.2.0",
		SubComponent: ".1.3.6.1.4.1.56.12.9.3.0",
	}
}

func Run(myConfigFromMain config.Config, alertsChannel chan types.Alert, waitGroup *sync.WaitGroup) {

	log.WithFields(logrus.Fields{"address": myConfigFromMain.SNMPTrapAddress}).Info("Starting the SNMP trapper")

	// Populate the config:
	myConfig = myConfigFromMain

	// Set up a channel to handle shutdown:
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)

	// Handle incoming alerts:
	go func() {
		for {
			select {

			case alert := <-alertsChannel:

				// Send a trap based on this alert:
				log.WithFields(logrus.Fields{"status": alert.Status}).Debug("Received an alert")
				sendTrap(alert)
			}
		}
	}()

	// Wait for shutdown:
	for {
		select {
		case <-signals:
			log.Warn("Shutting down the SNMP trapper")

			// Tell main() that we're done:
			waitGroup.Done()
			return
		}
	}

}
