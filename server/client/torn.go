package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/panaka13/torntools_server/server/model"
)

func queryTorn(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	ret := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func GetLogInfo(api string, user, from, to int32, logTypes []int) ([]model.TornLog, error) {
	url, _ := url.Parse("https://api.torn.com/user" + string(user))
	logTypeStrs := make([]string, 0, len(logTypes))

	query := url.Query()
	query.Set("selections", "log")
	query.Set("key", api)
	query.Set("from", string(from))
	query.Set("to", string(to))
	query.Set("log", strings.Join(logTypeStrs, ","))

	url.RawQuery = query.Encode()

	fmt.Println(url.String())
	resp, err := queryTorn(url.String())
	if err != nil {
		return nil, err
	}
	return handleLogResponse(resp)
}

func GetEventInfo(api string, user, from, to int32) ([]model.TornEvent, error) {
url, _ := url.Parse("https://api.torn.com/user" + string(user))
	query := url.Query()
	query.Set("selections", "log")
	query.Set("key", api)
	query.Set("from", string(from))
	query.Set("to", string(to))

	url.RawQuery = query.Encode()

	fmt.Println(url.String())
	resp, err := queryTorn(url.String())
	if err != nil {
		return nil, err
	}

	return handleEventResponse(resp)
}

func handleLogResponse(resp []byte) ([]model.TornLog, error) {
	var respMap map[string]json.RawMessage
	err := json.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
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