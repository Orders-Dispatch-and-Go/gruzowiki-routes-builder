package db

type PgrGeoJson struct {
	Type string `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}