package controller

import (
	"context"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	proto "github.com/panaka13/torntools_server/gen/torntools_proto"
	"github.com/panaka13/torntools_server/server/client"
	"github.com/panaka13/torntools_server/server/model"
	"github.com/panaka13/torntools_server/server/utility"
)

const (
	threeWayOrdinary = "3-Way Ordinary time"
	twoWayFull       = "2-Way Full event"
	asianHandicap    = "Asian Handicap"
	overUnder        = "Over/Under"
	btts             = "Both Teams to Score"
)

var (
	wonRegex          = regexp.MustCompile("^You won .* on your .* bet on")
	loseRegex         = regexp.MustCompile("^You lost your .* bet on")
	drawRegex         = regexp.MustCompile("^Your .* bet on .* was refunded$")
	winAmountRegex    = regexp.MustCompile(" on your \\$[0-9]{1,3}(,[0-9]{3})*")
	loseAmountRegex   = regexp.MustCompile("^You lost your \\$[0-9]{1,3}(,[0-9]{3})*")
	refundAmountRegex = regexp.MustCompile("^Your \\$[0-9]{1,3}(,[0-9]{3})*")
	descriptionRegex  = regexp.MustCompile(">.*</a>")
)

func getBookieResultFromTornLog(tornLog model.TornLog) proto.BookieResult {
	if strings.Contains(tornLog.Title, "Bookie win") {
		return proto.BookieResult_WIN
	}
	if strings.Contains(tornLog.Title, "Bookie refund") {
		return proto.BookieResult_REFUND
	}
	if strings.Contains(tornLog.Title, "Bookie lose") {
		return proto.BookieResult_LOSE
	}
	return proto.BookieResult_UNDECIDED
}

func getBookieResultFromTornEvent(tornEvent model.TornEvent) proto.BookieResult {
	if wonRegex.MatchString(tornEvent.Event) {
		return proto.BookieResult_WIN
	}
	if loseRegex.MatchString(tornEvent.Event) {
		return proto.BookieResult_LOSE
	}
	if drawRegex.MatchString(tornEvent.Event) {
		return proto.BookieResult_REFUND
	}
	return proto.BookieResult_UNDECIDED
}

func getBookieAmountFromTornEvent(tornEvent model.TornEvent, result proto.BookieResult) int32 {
	var str string
	switch result {
	case proto.BookieResult_WIN:
		match := winAmountRegex.FindString(tornEvent.Event)
		if len(match) == 0 {
			return 0
		}
		str = match[10:]
	case proto.BookieResult_LOSE:
		match := loseAmountRegex.FindString(tornEvent.Event)
		if len(match) == 0 {
			return 0
		}
		str = match[15:]
	case proto.BookieResult_REFUND:
		match := refundAmountRegex.FindString(tornEvent.Event)
		if len(match) == 0 {
			return 0
		}
		str = match[6:]
	default:
		str = "0"
	}
	value, err := strconv.ParseInt(strings.ReplaceAll(str, ",", ""), 10, 32)
	if err != nil {
		return 0
	}
	return int32(value)
}

func getBookieDescriptionFromTornEvent(tornEvent model.TornEvent) string {
	match := descriptionRegex.FindString(tornEvent.Event)
	if len(match) == 0 {
		return ""
	}
	return match[1 : len(match)-4]
}

func getBookieTypeFromTornEvent(tornEvent model.TornEvent) proto.BookieType {
	if strings.Contains(tornEvent.Event, threeWayOrdinary) {
		return proto.BookieType_THREE_WAY_ORDINARY
	}
	if strings.Contains(tornEvent.Event, twoWayFull) {
		return proto.BookieType_TWO_WAY_FULL
	}
	if strings.Contains(tornEvent.Event, asianHandicap) {
		return proto.BookieType_ASIAN_HANDICAP
	}
	if strings.Contains(tornEvent.Event, overUnder) {
		return proto.BookieType_OVER_UNDER
	}
	if strings.Contains(tornEvent.Event, btts) {
		return proto.BookieType_BTTS
	}
	return proto.BookieType_UNKNOWN
}

func tryMergeBookieResult(bookie *proto.Bookie, tornEvent model.TornEvent) bool {
	if bookie.Timestamp != tornEvent.Timestamp {
		return false
	}
	if bookie.Result != getBookieResultFromTornEvent(tornEvent) {
		return false
	}
	if bookie.Bet != getBookieAmountFromTornEvent(tornEvent, bookie.Result) {
		return false
	}
	bookie.Description = getBookieDescriptionFromTornEvent(tornEvent)
	bookie.BookieType = getBookieTypeFromTornEvent(tornEvent)
	return true
}

func mergeBookieResults(tornLogs []model.TornLog, tornEvents []model.TornEvent) []*proto.Bookie {
	sort.Slice(tornLogs, func(i, j int) bool {
		return tornLogs[i].Timestamp < tornLogs[j].Timestamp
	})
	bookies := make([]*proto.Bookie, 0, len(tornLogs))
	for _, tornLog := range tornLogs {
		bet, ok := tornLog.Data["bet"].(float64)
		if !ok {
			continue
		}
		oddsStr, ok := tornLog.Data["odds"].(string)
		if !ok {
			continue
		}
		odds, err := strconv.ParseFloat(oddsStr, 32)
		if err != nil {
			continue
		}
		var bookie proto.Bookie
		bookie.Id = tornLog.Id
		bookie.Timestamp = tornLog.Timestamp
		bookie.Bet = int32(bet)
		bookie.Odds = float32(odds)
		bookie.Result = getBookieResultFromTornLog(tornLog)
		eventFound := false
		for _, tornEvent := range tornEvents {
			if tryMergeBookieResult(&bookie, tornEvent) {
				eventFound = true
				break
			}
		}
		if eventFound {
			bookies = append(bookies, &bookie)
		}
	}
	return bookies
}

func (s TornToolServer) ViewBookieResult(ctx context.Context, request *proto.ViewBookieResultResquest) (*proto.ViewBookieResultResponse, error) {
	var response proto.ViewBookieResultResponse
	from := request.From.AsTime().Truncate(utility.DayDuration)
	to := request.To.AsTime().Truncate(utility.DayDuration).Add(utility.DayDuration - 1)

	tornLogs, err := client.GetLogInfo(request.Api, request.User, int32(from.Unix()), int32(to.Unix()), []int{8461, 8462, 8463})
	if err != nil {
		log.Printf("[ViewBookieResult] GetLogInfo error: %s\n", err.Error())
		response.Status = proto.Status_INTERNAL_ERROR
		response.Error = err.Error()
		return &response, nil
	}

	tornEvents, err := client.GetEventInfo(request.Api, request.User, int32(from.Unix()), int32(to.Unix()))
	if err != nil {
		response.Status = proto.Status_INTERNAL_ERROR
		response.Error = err.Error()
		return &response, nil
	}

	response.Results = mergeBookieResults(tornLogs, tornEvents)

	return &response, nil
}
