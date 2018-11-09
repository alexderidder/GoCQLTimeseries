package model

import (
	"encoding/json"
	"github.com/gocql/gocql"
)

type JSONUUID struct {
	Value gocql.UUID
	Valid bool
	Set   bool
}

func (i *JSONUUID) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Set = true

	if string(data) == "null" {
		// The key was set to null
		i.Valid = false
		return nil
	}

	// The key isn't set to null
	var temp gocql.UUID
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}
