package webhook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	logrus "github.com/Sirupsen/logrus"
	template "github.com/prometheus/alertmanager/template"
)

type WebhookHandler struct {
	AlertsChannel chan template.Alert
}

func (webhookHandler *WebhookHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// Read the request body:
	requestBody, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to read the request body")
		http.Error(responseWriter, "Failed to read the request body", http.StatusBadRequest)
		return
	}

	// Make a new alert:
	alert := &template.Alert{}

	// Unmarshal the request body into the alert:
	err = json.Unmarshal(requestBody, alert)
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Error("Failed to unmarshal the request body into an alert")
		http.Error(responseWriter, "Failed to unmarshal the request body into an alert", http.StatusBadRequest)
		return
	}

	// Report:
	log.WithFields(logrus.Fields{"requestbody": string(requestBody)}).Info("Received a valid webhook alert")

	// put the message onto the channel:
	webhookHandler.AlertsChannel <- *alert

}
