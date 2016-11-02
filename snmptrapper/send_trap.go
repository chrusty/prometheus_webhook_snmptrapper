package snmptrapper

import (
	types "github.com/chrusty/prometheus_webhook_snmptrapper/types"

	logrus "github.com/Sirupsen/logrus"
	snmpgo "github.com/k-sone/snmpgo"
)

func sendTrap(alert types.Alert) {

	// Prepare an SNMP handler:
	snmp, err := snmpgo.NewSNMP(snmpgo.SNMPArguments{
		Version:   snmpgo.V2c,
		Address:   myConfig.SNMPTrapAddress,
		Retries:   myConfig.SNMPRetries,
		Community: myConfig.SNMPCommunity,
	})
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to create snmpgo.SNMP object")
		return
	} else {
		log.WithFields(logrus.Fields{"address": myConfig.SNMPTrapAddress, "retries": myConfig.SNMPRetries, "community": myConfig.SNMPCommunity}).Debug("Created snmpgo.SNMP object")
	}

	// Build VarBind list:
	var varBinds snmpgo.VarBinds

	// System Uptime (ideally this would be the age of the alert):
	varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSysUpTime, snmpgo.NewTimeTicks(1000)))

	// The "enterprise OID" for the trap (rising/firing or falling/recovery):
	varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, trapOIDs.FiringTrap))

	// Figure out which "severity" value to send:
	switch {
	case alert.Status == "recovery":
		// "Any existing alerts with a matching Node, AlertGroup, AlertKey will be cleared":
		severity = 0
	case alert.Labels["severity"] == "info", alert.Labels["severity"] == "warning":
		// "The alert may provide useful information when attempting to determine the root cause of an issue":
		severity = 2
	case alert.Labels["severity"] == "minor":
		// "The alert represents a low priority issue that can be resolved during normal working hours":
		severity = 3
	case alert.Labels["severity"] == "major":
		// "An issue has occurred that is resilience affecting and requires immediate investigation":
		severity = 4
	case alert.Labels["severity"] == "critical":
		// "An issue has occurred that is service affecting and requires immediate investigation":
		severity = 5
	default:
		// "The severity of this alert is yet to be determined":
		severity = 1
	}

	// Insert the AlertManager variables:
	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Description, snmpgo.NewOctetString([]byte(alert.Annotations["description"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Severity, snmpgo.NewOctetString([]byte(alert.Labels["severity"]))))
	varBinds = append(varBinds, snmpgo.NewVarBind(trapOIDs.Service, snmpgo.NewOctetString([]byte(alert.Labels["service"]))))

	// Create an SNMP "connection":
	if err = snmp.Open(); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to open SNMP connection")
		return
	}
	defer snmp.Close()

	// Send the trap:
	if err = snmp.V2Trap(varBinds); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to send SNMP trap")
		return
	} else {
		log.WithFields(logrus.Fields{"status": alert.Status}).Info("It's a trap!")
	}
}
