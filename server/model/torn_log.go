package model

import (
	"encoding/json"
	"errors"
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
	if len(str) <= 2 {
		return  make([]TornLog, 0), nil
	}
	var tornLogMap map[string]TornLog
	err := json.Unmarshal(str, &tornLogMap)
	if err != nil {
		errorMsg := "str: " + string(str) + ", " + err.Error()
		return nil, errors.New(errorMsg)
	}
	tornLogs := make([]TornLog, 0, len(tornLogMap))
	for key, tornLog := range tornLogMap {
		tornLog.Id = key
		tornLogs = append(tornLogs, tornLog)
	}
	return tornLogs, nil
}
