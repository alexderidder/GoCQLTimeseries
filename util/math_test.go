package util

import (
	"GoCQLTimeSeries/datatypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinearInterpolation(t *testing.T) {
	solution, _ := LinearInterpolation(0,10,0,10,5)
	assert.Equal(t, float64(5) , solution , "Check linear interpolation Quadrant 1  (0,0) (10,10) 5 = 5")

	solution, _ =  LinearInterpolation(-10, 0,0,10,-5)
	assert.Equal(t, float64(5), solution, "Check linear interpolation Quadrant 2 (-10,0) (0,10) -5 = 5")

		solution, _ =LinearInterpolation(0,10,-10,0,5)
		assert.Equal(t, float64(-5), solution , "Check linear interpolation Quadrant 4  (0,-10) (10,0) 5 = -5")

	solution, _ = LinearInterpolation(-10,0,-10,0,-5)
		assert.Equal(t, float64(-5) , solution, "Check linear interpolation Quadrant 3 (-10,-10) (0,0) -5 = -5")

	solution, _ =LinearInterpolation(0,0,0,10,0)
	assert.Equal(t, float64(10), solution, "Check linear interpolation  (0,0) (0,10) 0 = 10")


	solution, _ = LinearInterpolation(0,10,0,10,10)
	assert.Equal(t, float64(10),solution, "Check linear interpolation  (0,0) (10,10) 10 = 10")

	_, err := LinearInterpolation(0, 10, 0, 10, 15)
	assert.Equal(t, datatypes.OutOfBounds, err, "Check linear interpolation  (0,0) (10,10) 15 = OutofBounds error")

	_, err = LinearInterpolation(0, 10, 0, 10, -5)
	assert.Equal(t, datatypes.OutOfBounds, err, "Check linear interpolation  (0,0) (10,10) -5 = OutofBounds error")


	maxFloat := 1.7976931348623157e+308
	solution, err = LinearInterpolation(0, 10, maxFloat, maxFloat, 5)
	assert.Equal(t, datatypes.Overflow, err, "Check linear interpolation  (0,max) (10,max) 5 = Overflow error")

	_, err = LinearInterpolation(10, 0, 0, 10, 5)
	assert.Equal(t, datatypes.WrongOrderXValues, err, "Check linear interpolation  (10,0) (0,10) 5 = x1SmallerThenx1 error")

}
