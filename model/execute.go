package model

import "GoCQLTimeSeries/datatypes"

type Execute interface {
	Execute() ([]byte, datatypes.Error)
}