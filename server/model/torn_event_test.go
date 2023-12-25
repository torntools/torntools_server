package model

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func isTornEventEqual(this, other TornEvent) bool {
	if this.Id != other.Id {
		return false
	}
	if this.Timestamp != other.Timestamp{
		return false
	}
	if this.Event != other.Event {
		return false
	}
	return true
}

func TestParseTornEvent(t *testing.T) {
	data := []byte(`
{
    "GwZ3pxcRwiFoeyC86vlA": {
      "timestamp": 1703478314,
      "event": "You have been given 10x Snowballs and 3x Special Refills. Have a wonderful day!"
    },
    "iYEqxSgkKpA4B0FHGkJE": {
      "timestamp": 1703436553,
      "event": "You're <a href = http://www.torn.com/\"http://www.torn.com/http://www.torn.com/profiles.php?XID=7045\">GanjaMama</a>'s Secret Santa. Are you going to make their holiday special, or will you be the grinch?"
    },
    "oUHbIURCzmGgaTASnDN6": {
      "timestamp": 1703433528,
      "event": "You're <a href = http://www.torn.com/\"http://www.torn.com/http://www.torn.com/profiles.php?XID=7045\">GanjaMama</a>'s Secret Santa, oh what fun, will you send a gift or make them run?"
    }
}
`)

	expected := []TornEvent{
		{
			Id: "GwZ3pxcRwiFoeyC86vlA",
			Timestamp: 1703478314,
			Event: "You have been given 10x Snowballs and 3x Special Refills. Have a wonderful day!",
		},
		{
			Id: "iYEqxSgkKpA4B0FHGkJE",
			Timestamp: 1703436553,
			Event: "You're <a href = http://www.torn.com/\"http://www.torn.com/http://www.torn.com/profiles.php?XID=7045\">GanjaMama</a>'s Secret Santa. Are you going to make their holiday special, or will you be the grinch?",
		},
		{
			Id: "oUHbIURCzmGgaTASnDN6",
			Timestamp: 1703433528,
			Event: "You're <a href = http://www.torn.com/\"http://www.torn.com/http://www.torn.com/profiles.php?XID=7045\">GanjaMama</a>'s Secret Santa, oh what fun, will you send a gift or make them run?",
		},
	}

	tornEvents, err := ParseTornEvents(data)
	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, len(tornEvents), len(expected), "Expect 3 log entries")
	for i := 0; i < len(tornEvents); i++ {
		foundEqual := false
		for j := 0; j < len(expected); j++ {
			if isTornEventEqual(tornEvents[i], expected[j]) {
				foundEqual = true
				break
			}
		}
		assert.True(t, foundEqual, "Parsed event not equals expected value")
	}
}
