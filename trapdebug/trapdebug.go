package main

import (
	"flag"
	"net"

	logrus "github.com/sirupsen/logrus"
	gosnmptrap "github.com/ebookbug/gosnmptrap"
)

var (
	listenPort = flag.Int("listenport", 162, "Port to listen for traps on")
)

func init() {

	// Process the command-line parameters:
	flag.Parse()

	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {

	logrus.WithFields(logrus.Fields{"port": *listenPort}).Info("Starting SNMP TrapDebugger")

	// Open a UDP socket:
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: *listenPort,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Fatal("Error opening socket")
	}
	defer socket.Close()

	// Loop forever:
	for {
		// Make a buffer to read into:
		buf := make([]byte, 2048)

		// Read from the socket:
		read, from, _ := socket.ReadFromUDP(buf)

		// Report that we have data:
		logrus.WithFields(logrus.Fields{"client": from.IP}).Debug("Data received")

		// Handle the data:
		go HandleUdp(buf[:read])
	}
}

// Handle SNMP data:
func HandleUdp(data []byte) {

	// Attempt to parse the SNMP data:
	trap, err := gosnmptrap.ParseUdp(data)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Unable to parse SNMP data")
		return
	}

	// Dump the metadata:
	logrus.WithFields(logrus.Fields{"version": trap.Version, "community": trap.Community, "enterprise_id": trap.EnterpriseId, "address": trap.Address}).Info("SNMP trap received")
	logrus.WithFields(logrus.Fields{"general": trap.GeneralTrap, "special": trap.SpeicalTrap}).Info("SNMP trap values")

	// Dump the values:
	for trapOID, trapValue := range trap.Values {
		logrus.WithFields(logrus.Fields{"OID": trapOID, "value": trapValue}).Info("Trap variable")
	}
}
