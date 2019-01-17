package util

func LinearInterpolation(x0 int64, x1 int64, y0 float64, y1 float64, x int64) (float64){
	return (y0*float64(x1-x) + y1*float64(x-x0)) / float64(x1-x0)
}
