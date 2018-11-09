package model

import (
	"time"
)

type DeleteJSON struct {
	StoneID   JSONUUID `json:"stoneIDs"`
	Types     []string   `json:"types"`
	StartTime time.Time  `json:"startTime"`
	EndTime   time.Time  `json:"endTime"`
}
