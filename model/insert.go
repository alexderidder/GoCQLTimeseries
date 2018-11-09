package model

import (
	"time"
)

type InsertJSON struct {
	StoneID JSONUUID `json:"stoneID"`
	Data    []struct {
		Time        time.Time   `json:"time"`
		Watt        JSONFloat32 `json:"watt"`
		PowerFactor JSONFloat32 `json:"pf"`
		KWH         JSONFloat32 `json:"kWh"`
	} `json:"data"`
}


