package webhook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	invalidPayload = `cruft goes here`
	validPayload   = `
{
  "version": "2",
  "status": "firing",
  "alerts": [
    {
      "labels": {
        "datacentre": "dc1"
      },
      "annotations": {
        "summary": "this is a test cruft alert"
      },
      "startsAt": "2016-10-27T14:27:00Z",
      "endsAt": "2016-10-27T14:27:00Z"
    },
    {
      "labels": {
        "datacentre": "dc1"
      },
      "annotations": {
        "summary": "this is a test garbage alert"
      },
      "startsAt": "2016-10-27T14:27:00Z",
      "endsAt": "2016-10-27T14:27:00Z"
    }
  ]
}
`
)

func TestWebhook(t *testing.T) {

	testWebhookValidPayload(t)
	testWebhookInvalidPayload(t)

}

func testWebhookValidPayload(t *testing.T) {

	// Validate the payload:
	err, alerts := validatePayload([]byte(validPayload))

	// Assess the results:
	assert.NoError(t, err, "Unable to validate a valid payload")
	assert.EqualValues(t, 2, len(alerts), "Got the wrong number of validated alerts")
}

func testWebhookInvalidPayload(t *testing.T) {

	// Validate the payload:
	err, alerts := validatePayload([]byte(invalidPayload))

	// Assess the results:
	assert.Error(t, err, "Validated an invalid JSON payload")
	assert.EqualValues(t, 0, len(alerts), "Got the wrong number of validated alerts")
}
