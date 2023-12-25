package model

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func isTornLogEqual(this, other TornLog) bool {
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

func TestParseTornLog(t *testing.T) {
	data := []byte(`
{
    "jlEi5fKSolk8YGo5U4Kf": {
        "log": 9010,
        "title": "Crime success",
        "timestamp": 1703421963,
        "category": "Crimes",
        "data": {
            "crime_action": "brute forcing a password",
            "outcome": 14940,
            "nerve": 7,
            "unique": ""
        },
        "params": {
            "color": "green"
        }
    },
    "uEV9z2CSIqWbDZtY7H1h": {
        "log": 9010,
        "title": "Crime success",
        "timestamp": 1703421961,
        "category": "Crimes",
        "data": {
            "crime_action": "brute forcing a password",
            "outcome": 14937,
            "nerve": 7,
            "unique": ""
        },
        "params": {
            "color": "green"
        }
    },
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
}`)

	expected := []TornLog{
		{
			Id:        "jlEi5fKSolk8YGo5U4Kf",
			Log:       9010,
			Category:  "Crimes",
			Title:     "Crime success",
			Timestamp: 1703421963,
			Data: map[string]interface{}{
				"crime_action": "brute forcing a password",
				"outcome":      14940,
				"nerve":        7,
				"unique":       "",
			},
			Params: map[string]interface{}{
				"color": "green",
			},
		},
		{
			Id:        "uEV9z2CSIqWbDZtY7H1h",
			Log:       9010,
			Category:  "Crimes",
			Title:     "Crime success",
			Timestamp: 1703421961,
			Data: map[string]interface{}{
				"crime_action": "brute forcing a password",
				"outcome":      14937,
				"nerve":        7,
				"unique":       "",
			},
			Params: map[string]interface{}{
				"color": "green",
			},
		},
		{
			Id:        "77BD5LZvuvTU04LKqQI1",
			Log:       7205,
			Category:  "Messages",
			Title:     "Message read",
			Timestamp: 1703417604,
			Data: map[string]interface{}{
				"message": 245501027,
			},
			Params: map[string]interface{}{},
		},
	}

	tornLogs, err := ParseTornLogs(data)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, len(tornLogs), len(expected), "Expect 3 log entries")
	for i := 0; i < len(tornLogs); i++ {
		foundEqual := false
		for j := 0; j < len(expected); j++ {
			if isTornLogEqual(tornLogs[i], expected[j]) {
				foundEqual = true
				break
			}
		}
		assert.True(t, foundEqual, "Parsed log not equals expected value")
	}
}
