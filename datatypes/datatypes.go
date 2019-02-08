package datatypes

type Data struct {
	Time  int64 `json:"time"`
	Value struct {
		Wattage     float32 `json:"w,omitempty"`
		PowerFactor float32 `json:"pf,omitempty"`
		KWH         float64 `json:"kWh,omitempty"`
	} `json:"value"`
}