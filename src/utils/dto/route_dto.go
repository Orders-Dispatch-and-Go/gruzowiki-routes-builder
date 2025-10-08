package dto

type RouteDTO struct {
	Way [][]float64 `json:"way"`
	Graph []int32 `json:"graph"`
}

func (routeDTO *RouteDTO) IsDtoValid() (bool) {
	return routeDTO == nil || routeDTO.Graph != nil && routeDTO.Way != nil;
}