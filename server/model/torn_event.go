package model

import "encoding/json"

type TornEvent struct {
	Id string
	Timestamp int32
	Event string
}

func ParseTornEvents(str []byte) ([]TornEvent, error) {
	var tornLogMap map[string]TornEvent
	err := json.Unmarshal(str, &tornLogMap)
	if err != nil {
		return nil, err
	}
	tornLogs := make([]TornEvent, 0, len(tornLogMap))
	for key, tornLog := range tornLogMap {
		tornLog.Id = key
		tornLogs = append(tornLogs, tornLog)
	}
	return tornLogs, nil
}