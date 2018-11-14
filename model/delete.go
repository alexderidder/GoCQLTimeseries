package model

import (
	"time"
)

type DeleteJSON struct {
	StoneID   JSONUUID `json:"stoneID"`
	Types     []string   `json:"types"`
	StartTime time.Time  `json:"startTime"`
	EndTime   time.Time  `json:"endTime"`
}
