package dto

type CreateRouteDTO struct {
	Stops [][]float64 `json:"stops"`
}