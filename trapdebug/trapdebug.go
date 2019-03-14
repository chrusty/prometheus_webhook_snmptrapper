package main

import (
	"flag"
	"net"
	"strconv"
	"strings"
	"os"

	logrus "github.com/Sirupsen/logrus"
	gosnmptrap "github.com/ebookbug/gosnmptrap"
)

var (
	listenPort = flag.Int("listenport", 162, "Port to listen for traps on")
	listenType = flag.String("listentype", "udp", "Network udp|udp4|udp6|tcp|tcp4|tcp6")
	log        = logrus.WithFields(logrus.Fields{"logger": "trapdebug-listener"})
)

func init() {

	// Process the command-line parameters:
	flag.Parse()

	// Set the log-level:
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {

    *listenType = strings.ToLower(*listenType)

    // UDP
	if (*listenType == "udp" || *listenType == "udp4" || *listenType == "udp6") {

	  // Open a UDP listener:
	  log.WithFields(logrus.Fields{"Protocol": *listenType, "Port": *listenPort}).Info("Open a UDP listener")

      udpAddr, err := net.ResolveUDPAddr(*listenType, ":"+strconv.Itoa(*listenPort))
      if err != nil {
		  logrus.WithFields(logrus.Fields{"error": err}).Fatal("Error creating UDPAddr")
	  }

  	  listen, err := net.ListenUDP(*listenType, udpAddr)
	  if err != nil {
		  logrus.WithFields(logrus.Fields{"error": err}).Fatal("Error opening UDP listener")
	  }

	  defer listen.Close()
	  log.WithFields(logrus.Fields{"network": *listenType, "port": *listenPort}).Info("UDP listener started")

	  // Loop forever:
	  for {
		  // Make a buffer to read into:
		  buf := make([]byte, 2048)

		  // Read from the listener
		  read, from, _ := listen.ReadFromUDP(buf)

		  // Report that we have data:
		  logrus.WithFields(logrus.Fields{"client": from.IP}).Debug("Data received")

		  // Handle the data:
		  go HandleSNMPdata(buf[:read])
	  }
	}

    // TCP
	if (*listenType == "tcp" || *listenType == "tcp4" || *listenType == "tcp6") {

	  // Open a TCP listener:
	  log.WithFields(logrus.Fields{"Protocol": *listenType, "Port": *listenPort}).Info("Open a TCP listener")

      tcpAddr, err := net.ResolveTCPAddr(*listenType, ":"+strconv.Itoa(*listenPort))
      if err != nil {
		  logrus.WithFields(logrus.Fields{"error": err}).Fatal("Error creating TCPAddr")
	  }

  	  listen, err := net.ListenTCP(*listenType, tcpAddr)
	  if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("Open TCP listener failed")
		return
	  }

	  defer listen.Close()
	  log.WithFields(logrus.Fields{"network": *listenType, "port": *listenPort}).Info("TCP listener started")

      // loop forever
	  for {
		conn, err := listen.Accept()
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Error("TCP listener failed")
			continue
		}
		// Report that we have data:
		logrus.WithFields(logrus.Fields{"port": *listenPort}).Debug("Data received")

		go HandleTCPdata(conn)
	  }
	}
}

// handle TCP data
func HandleTCPdata(conn net.Conn) {

	defer conn.Close()

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
      log.WithFields(logrus.Fields{"Error": err}).Error("error reading from TCP connection")
    }
    go HandleSNMPdata(buf[:n])
}

// Handle SNMP data:
func HandleSNMPdata(data []byte) {

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
	go writeToFile(string(data[:len(data)]))
}

// write notification to file
func writeToFile(data string) {

    filename := "/tmp/snmptrap.txt"
    log.WithFields(logrus.Fields{"File": filename}).Info("creating file")

	f, err := os.Create(filename)
    if err != nil {
      log.WithFields(logrus.Fields{"Error": err}).Error("error creating file")
      return
    }

    l, err := f.WriteString(strconv.Quote(data))
    if err != nil {
      log.WithFields(logrus.Fields{"Error": err}).Error("error writing data")
      f.Close()
      return
    }

	log.WithFields(logrus.Fields{"Length": l}).Info("bytes written successfully")
    err = f.Close()
    if err != nil {
      log.WithFields(logrus.Fields{"Error": err}).Error("error closing file")
      return
	}
}
