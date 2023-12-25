package model

import (
	"encoding/json"
)

type TornLog struct {
	Id        string
	Log       int16
	Category  string
	Title     string
	Timestamp int32
	Data      map[string]interface{}
	Params    map[string]interface{}
}

func ParseTornLogs(str []byte) ([]TornLog, error) {
	var tornLogMap map[string]TornLog
	err := json.Unmarshal(str, &tornLogMap)
	if err != nil {
		return nil, err
	}
	tornLogs := make([]TornLog, 0, len(tornLogMap))
	for key, tornLog := range tornLogMap {
		tornLog.Id = key
		tornLogs = append(tornLogs, tornLog)
	}
	return tornLogs, nil
}
