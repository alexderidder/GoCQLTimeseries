package util

import (
	"GoCQLTimeSeries/datatypes"
	"math"
)

func LinearInterpolation(x1 int64, x2 int64, y1 float64, y2 float64, x int64) (float64, datatypes.Error) {
	if y1 > math.MaxFloat64 / 2 || y2 > math.MaxFloat64 / 2 {
		return 0, datatypes.Overflow
	}
	if x1 > x2 {
		return 0, datatypes.WrongOrderXValues
	}
	if x < x1 || x > x2 {
		return 0, datatypes.OutOfBounds
	}

	if x1 == x2 {
		return y2, datatypes.NoError
	}
	return (y1*float64(x2-x) + y2*float64(x-x1)) / float64(x2-x1), datatypes.NoError
}
