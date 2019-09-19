package snmptrapper

import (
	"os"
	"os/signal"
	"sync"

	config "github.com/chrusty/prometheus_webhook_snmptrapper/config"
	types "github.com/chrusty/prometheus_webhook_snmptrapper/types"

	logrus "github.com/sirupsen/logrus"
	snmpgo "github.com/k-sone/snmpgo"
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
	trapOIDs.FiringTrap, _ = snmpgo.NewOid("1.3.6.1.3.1977.1.0.1")
	trapOIDs.RecoveryTrap, _ = snmpgo.NewOid("1.3.6.1.3.1977.1.0.2")
	trapOIDs.Instance, _ = snmpgo.NewOid("1.3.6.1.3.1977.1.1.1")
	trapOIDs.Service, _ = snmpgo.NewOid("1.3.6.1.3.1977.1.1.2")
	trapOIDs.Location, _ = snmpgo.NewOid("1.3.6.1.3.1977.1.1.3")
	trapOIDs.Severity, _ = snmpgo.NewOid("1.3.6.1.3.1977.1.1.4")
	trapOIDs.Description, _ = snmpgo.NewOid("1.3.6.1.3.1977.1.1.5")
	trapOIDs.JobName, _ = snmpgo.NewOid("1.3.6.1.3.1977.1.1.6")
	trapOIDs.TimeStamp, _ = snmpgo.NewOid("1.3.6.1.3.1977.1.1.7")
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
