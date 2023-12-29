package client

import (
	"strings"
	"testing"

	"github.com/panaka13/torntools_server/server/model"
	"github.com/stretchr/testify/assert"
)

func TestHandleLogErrorResponse(t *testing.T) {
	resp := []byte(`{"error":{"code":2,"error":"Incorrect key"}}`)
	_, err := handleLogResponse(resp)
	assert.Error(t, err, "Should return error")
	assert.True(t, strings.Contains(err.Error(), "Incorrect key"), `Error should be "Incorrect key"`)
}

func TestHandleLogResponse(t *testing.T) {
	resp := []byte(`
{
  "log": {
    "77BD5LZvuvTU04LKqQI1": {
      "log": 7205,
      "title": "Message read",
      "timestamp": 1703417604,
      "category": "Messages",
      "data": {
        "message": 245501027
      },
      "params": {}
    }
  }
}`)
	tornLogs, err := handleLogResponse(resp)
	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, len(tornLogs), 1, "Expected 1 instance of tornLog")
	expected := model.TornLog{
		Id: "77BD5LZvuvTU04LKqQI1",
		Log: 7205,
		Title: "Message read",
		Category: "Messages",
		Timestamp: 1703417604,
		Data: map[string]interface{}{"message": 245501027},
	}
	assert.True(t, isTornLogEqual(tornLogs[0], expected), "Parsed TornLog not match expected")
}

func TestHandleEventErrorResponse(t *testing.T) {
	resp := []byte(`{"error":{"code":2,"error":"Incorrect key"}}`)
	_, err := handleLogResponse(resp)
	assert.Error(t, err, "Should return error")
	assert.True(t, strings.Contains(err.Error(), "Incorrect key"), `Error should be "Incorrect key"`)
}

func TestHandleEventResponse(t *testing.T) {
	resp := []byte(`
{
  "events": {
    "GwZ3pxcRwiFoeyC86vlA": {
      "timestamp": 1703478314,
      "event": "You have been given 10x Snowballs and 3x Special Refills. Have a wonderful day!"
    }
  }
}
`)
	tornEvents, err := handleEventResponse(resp)
	assert.Nil(t, err, "Should not return error")
	assert.Equal(t, len(tornEvents), 1, "Expect 1 instance of tornEvent")
	expected := model.TornEvent{
		Id: "GwZ3pxcRwiFoeyC86vlA",
		Timestamp: 1703478314,
		Event: "You have been given 10x Snowballs and 3x Special Refills. Have a wonderful day!",
	}
	assert.True(t, isTornEventEqual(tornEvents[0], expected), "Parsed Event not match expected")
}

func isTornLogEqual(this, other model.TornLog) bool {
	if this.Id != other.Id {
		return false
	}
	if this.Log != other.Log {
		return false
	}
	if this.Category != other.Category {
		return false
	}
	if this.Title != other.Title {
		return false
	}
	if this.Timestamp != other.Timestamp {
		return false
	}
	return true
}

func isTornEventEqual(this, other model.TornEvent) bool {
	if this.Id !=  other.Id {
		return false
	}
	if this.Timestamp!=  other.Timestamp {
		return false
	}
	if this.Event != other.Event {
		return false
	}
	return true
}