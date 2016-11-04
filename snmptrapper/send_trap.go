package snmptrapper

import (
	"bytes"
	"fmt"
	"os/exec"

	types "github.com/chrusty/prometheus_webhook_snmptrapper/types"

	logrus "github.com/Sirupsen/logrus"
)

func sendTrap(alert types.Alert) {
	var genericTrap string = "6"
	var specificTrap string = "1"

	// Figure out which "severity" value to send:
	switch {
	case alert.Status == "recovery":
		// "Any existing alerts with a matching Node, AlertGroup, AlertKey will be cleared":
		specificTrap = "0"
	case alert.Labels["severity"] == "info", alert.Labels["severity"] == "warning":
		// "The alert may provide useful information when attempting to determine the root cause of an issue":
		specificTrap = "2"
	case alert.Labels["severity"] == "minor":
		// "specificTrapThe alert represents a low priority issue that can be resolved during normal working hours":
		specificTrap = "3"
	case alert.Labels["severity"] == "major":
		// "An issue has occurred that is resilience affecting and requires immediate investigation":
		specificTrap = "4"
	case alert.Labels["severity"] == "critical":
		// "An issue has occurred that is service affecting and requires immediate investigation":
		specificTrap = "5"
	default:
		// "The severity of this alert is yet to be determined":
		specificTrap = "1"
	}

	// Prepare to send the TRAP using Net-SNMP's "snmptrap" command (because I can't find a library capable of sending V1 traps in GoLang):
	var stdout, stderr bytes.Buffer
	var arguments = make(map[string]string)
	arguments["snmpTrapVersion"] = "-v1"
	arguments["snmpCommunity"] = fmt.Sprintf("-c %v", myConfig.SNMPCommunity)
	arguments["snmpTrapdAddress"] = myConfig.SNMPTrapAddress
	arguments["snmpTrapdOID"] = trapOIDs.TrapOID
	arguments["agentAddress"] = "127.0.0.1" // alert.Address
	arguments["genericTrap"] = genericTrap
	arguments["specificTrap"] = specificTrap
	arguments["uptime"] = "0"
	arguments["oidComponent"] = trapOIDs.Component
	arguments["oidComponentType"] = "s"
	arguments["oidComponentValue"] = fmt.Sprintf("'%v'", alert.Labels["instance"])
	arguments["oidMessage"] = trapOIDs.Message
	arguments["oidMessageType"] = "s"
	arguments["oidMessageValue"] = fmt.Sprintf("'%v'", alert.Annotations["description"])
	arguments["oidSubComponent"] = trapOIDs.SubComponent
	arguments["oidSubComponentType"] = "s"
	arguments["oidSubComponentValue"] = fmt.Sprintf("'%v'", alert.Labels["service"])

	// Trap command:
	netSNMPTrapCommand := exec.Command(
		myConfig.SNMPTrapBinary,
		arguments["snmpTrapVersion"],
		arguments["snmpCommunity"],
		arguments["snmpTrapdAddress"],
		arguments["snmpTrapdOID"],
		arguments["agentAddress"],
		arguments["genericTrap"],
		arguments["specificTrap"],
		arguments["uptime"],
		arguments["oidComponent"],
		arguments["oidComponentType"],
		arguments["oidComponentValue"],
		arguments["oidMessage"],
		arguments["oidMessageType"],
		arguments["oidMessageValue"],
		arguments["oidSubComponent"],
		arguments["oidSubComponentType"],
		arguments["oidSubComponentValue"],
	)
	netSNMPTrapCommand.Stdout = &stdout
	netSNMPTrapCommand.Stderr = &stderr

	// Send the trap:
	err := netSNMPTrapCommand.Run()
	if err != nil {
		log.WithFields(logrus.Fields{"error": err, "stdout": stdout.String(), "stderr": stderr.String(), "command": netSNMPTrapCommand.Path, "args": netSNMPTrapCommand.Args}).Error("Failed to send SNMP trap")
		return
	} else {
		log.WithFields(logrus.Fields{"status": alert.Status, "specific_trap": specificTrap, "generic_trap": genericTrap, "severity": alert.Labels["severity"]}).Info("It's a trap!")
	}
}
