package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/panaka13/torntools_server/server/model"
)

func queryTorn(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	ret, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func GetLogInfo(api string, user, from, to int32, logTypes []int) ([]model.TornLog, error) {
	queryUrl, _ := url.Parse("https://api.torn.com/user/" + strconv.Itoa(int(user)))
	logTypeStrs := make([]string, 0, len(logTypes))
	for _, log := range logTypes {
		logTypeStrs = append(logTypeStrs, strconv.Itoa(int(log)))
	}
	query := queryUrl.Query()
	query.Set("selections", "log")
	query.Set("key", api)
	query.Set("from", strconv.Itoa(int(from)))
	query.Set("to", strconv.Itoa(int(to)))
	query.Set("log", strings.Join(logTypeStrs, ","))

	queryUrl.RawQuery = query.Encode()

	resp, err := queryTorn(queryUrl.String())
	if err != nil {
		return nil, err
	}
	return handleLogResponse(resp)
}

func GetEventInfo(api string, user, from, to int32) ([]model.TornEvent, error) {
	queryUrl, _ := url.Parse("https://api.torn.com/user/" + strconv.Itoa(int(user)))
	query := queryUrl.Query()
	query.Set("selections", "events")
	query.Set("key", api)
	query.Set("from", strconv.Itoa(int(from)))
	query.Set("to", strconv.Itoa(int(to)))

	queryUrl.RawQuery = query.Encode()

	resp, err := queryTorn(queryUrl.String())
	if err != nil {
		return nil, err
	}

	return handleEventResponse(resp)
}

func handleLogResponse(resp []byte) ([]model.TornLog, error) {
	var respMap map[string]json.RawMessage
	err := json.Unmarshal(resp, &respMap)
	if err != nil {
		errorMsg := "response: " + string(resp) + err.Error()
		return nil, errors.New(errorMsg)
	}
	errorMsg, exist := respMap["error"]
	if exist {
		return nil, errors.New(string(errorMsg))
	}

	logMsg, exist := respMap["log"]
	if !exist {
		return nil, errors.New("Torn response not contain log")
	}
	return model.ParseTornLogs([]byte(logMsg))
}

func handleEventResponse(resp []byte) ([]model.TornEvent, error) {
	var respMap map[string]json.RawMessage
	err := json.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
	}
	errorMsg, exist := respMap["error"]
	if exist {
		return nil, errors.New(string(errorMsg))
	}

	eventMsg, exist := respMap["events"]
	if !exist {
		return nil, errors.New("Torn response not contain events")
	}
	return model.ParseTornEvents([]byte(eventMsg))
}

// export async function getEventInfo(api: string, user: number, from: number, to: number) {
//   const url = buildUrl("https://api.torn.com/", {
//     path: "user/" + String(user),
//     queryParams: {
//       selections: "events",
//       key: api,
//       from,
//       to,
//       comment: "panaka_13"
//     }
//   });
//   const [error, json] = await queryTorn(url);
//   if (error !== null) {
//     return [error, json];
//   }
//   return [null, json.events];
// }
