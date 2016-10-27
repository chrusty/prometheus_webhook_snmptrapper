package snmptrapper

import (
	logrus "github.com/Sirupsen/logrus"
	snmpgo "github.com/k-sone/snmpgo"
	template "github.com/prometheus/alertmanager/template"
)

func sendTrap(alert template.Alert) {
	snmp, err := snmpgo.NewSNMP(snmpgo.SNMPArguments{
		Version:   snmpgo.V2c,
		Address:   myConfig.SNMPTrapAddress,
		Retries:   myConfig.SNMPRetries,
		Community: myConfig.SNMPCommunity,
	})
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to create snmpgo.SNMP object")
		return
	}

	// Build VarBind list
	var varBinds snmpgo.VarBinds
	varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSysUpTime, snmpgo.NewTimeTicks(1000)))

	oid, _ := snmpgo.NewOid("1.3.6.1.6.3.1.1.5.3")
	varBinds = append(varBinds, snmpgo.NewVarBind(snmpgo.OidSnmpTrap, oid))

	oid, _ = snmpgo.NewOid("1.3.6.1.2.1.2.2.1.1.2")
	varBinds = append(varBinds, snmpgo.NewVarBind(oid, snmpgo.NewInteger(2)))

	oid, _ = snmpgo.NewOid("1.3.6.1.2.1.31.1.1.1.1.2")
	varBinds = append(varBinds, snmpgo.NewVarBind(oid, snmpgo.NewOctetString([]byte("eth0"))))

	if err = snmp.Open(); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to open SNMP connection")
		return
	}
	defer snmp.Close()

	if err = snmp.V2Trap(varBinds); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to send SNMP trap")
		return
	} else {
		log.Info("It's a trap!")
	}
}
