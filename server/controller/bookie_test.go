package controller

import (
	"testing"

	proto "github.com/panaka13/torntools_server/gen/torntools_proto"
	"github.com/panaka13/torntools_server/server/model"
	"github.com/stretchr/testify/assert"
)

const (
	wonEvent        = `You won $5,670,000 on your $3,500,000 Tottenham Hotspur (3-Way Ordinary time) bet on <a href = http://www.torn.com/\"http://www.torn.com/http://www.torn.com/page.php?sid=bookie#/your-bets/4193708\">Tottenham Hotspur v Everton</a>`
	wonBet          = 3500000
	wonDescription  = "Tottenham Hotspur v Everton"
	lostEvent       = `You lost your $3,000,000 Manchester United (3-Way Ordinary time) bet on <a href = http://www.torn.com/\"http://www.torn.com/http://www.torn.com/page.php?sid=bookie#/your-bets/4193709\">West Ham United v Manchester United</a>`
	lostBet         = 3000000
	lostDescription = "West Ham United v Manchester United"
	drawEvent       = `Your $5,000,000 Liverpool (-1) (Asian Handicap 1 Ordinary time) bet on <a href = http://www.torn.com/\"http://www.torn.com/http://www.torn.com/page.php?sid=bookie#/your-bets/4193683\">Crystal Palace v Liverpool</a> was refunded`
	drawBet         = 5000000
	drawDescription = "Crystal Palace v Liverpool"
	noBet           = 0
)

func TestGetBookieResultFromTornLog(t *testing.T) {
	type testCaseT struct {
		title  string
		expect proto.BookieResult
	}
	testCases := []testCaseT{
		{
			title:  "Bookie win",
			expect: proto.BookieResult_WIN,
		},
		{
			title:  "Bookie lose",
			expect: proto.BookieResult_LOSE,
		},
		{
			title:  "Bookie refund",
			expect: proto.BookieResult_REFUND,
		},
		{
			title:  "Message read",
			expect: proto.BookieResult_UNDECIDED,
		},
	}
	for _, tc := range testCases {
		var tornLog model.TornLog
		tornLog.Title = tc.title
		result := getBookieResultFromTornLog(tornLog)
		assert.Equal(t, tc.expect, result)
	}
}

func TestGetBookieResultFromTornEvent(t *testing.T) {
	type testCaseT struct {
		event  string
		expect proto.BookieResult
	}
	testCases := []testCaseT{
		{
			event:  wonEvent,
			expect: proto.BookieResult_WIN,
		},
		{
			event:  lostEvent,
			expect: proto.BookieResult_LOSE,
		},
		{
			event:  drawEvent,
			expect: proto.BookieResult_REFUND,
		},
	}
	for _, tc := range testCases {
		var tornEvent model.TornEvent
		tornEvent.Event = tc.event
		result := getBookieResultFromTornEvent(tornEvent)
		assert.Equal(t, tc.expect, result)
	}
}

func TestGetBookieAmountFromTornEvent(t *testing.T) {
	type testCaseT struct {
		event        string
		bookieResult proto.BookieResult
		expect       int32
	}
	testCases := []testCaseT{
		{
			event:        wonEvent,
			bookieResult: proto.BookieResult_WIN,
			expect:       wonBet,
		},
		{
			event:        lostEvent,
			bookieResult: proto.BookieResult_LOSE,
			expect:       lostBet,
		},
		{
			event:        drawEvent,
			bookieResult: proto.BookieResult_REFUND,
			expect:       drawBet,
		},
	}
	for _, tc := range testCases {
		var tornEvent model.TornEvent
		tornEvent.Event = tc.event
		result := getBookieAmountFromTornEvent(tornEvent, tc.bookieResult)
		assert.Equal(t, tc.expect, result)
	}
}

func TestGetBookieDescriptionFromTornEvent(t *testing.T) {
	type testCaseT struct {
		event  string
		expect string
	}
	testCases := []testCaseT{
		{
			event:  wonEvent,
			expect: wonDescription,
		},
		{
			event:  lostEvent,
			expect: lostDescription,
		},
		{
			event:  drawEvent,
			expect: drawDescription,
		},
	}
	for _, tc := range testCases {
		var tornEvent model.TornEvent
		tornEvent.Event = tc.event
		result := getBookieDescriptionFromTornEvent(tornEvent)
		assert.Equal(t, tc.expect, result)
	}
}
