package dto

type RouteDTO struct {
	Way [][]float64 `json:"way"`
	Distance float64 `json:"distance"` // в метрах
	Duration float64 `json:"duration"` // в секундах
}

func (routeDTO *RouteDTO) IsDtoValid() (bool) {
	return routeDTO != nil && routeDTO.Way != nil;
}