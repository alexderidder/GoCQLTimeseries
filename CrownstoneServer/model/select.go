package model

import (
	"github.com/gocql/gocql"
	"time"
)

type RequestSelectJSON struct {
	StoneIDs  []JSONUUID `json:"stoneIDs"`
	Types     []string     `json:"types"`
	StartTime time.Time    `json:"startTime"`
	EndTime   time.Time    `json:"endTime"`
	Interval  uint32       `json:"interval"`
}

type ResponseSelectJSON struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Interval  uint32    `json:"interval"`
	Stones    []Stone   `json:"stones"`
}

type Stone struct {
	StoneID gocql.UUID `json:"stoneID"`
	Fields  []Field    `json:"fields"`
}

type Field struct {
	Field string `json:"field"`
	Data  []Data `json:"Data"`
}

type Data struct {
	Time  time.Time `json:"time"`
	Value float32   `json:"value"`
}

const (
	UnitW      string = "w"
	Unitpf     string = "pf"
	UnitWAndpf string = "w_pf"
	UnitkWh    string = "kwh"
)
