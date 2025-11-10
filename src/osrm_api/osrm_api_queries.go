package osrmapi

import (
	"GruzowikiRoutesGenerator/utils/dto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type OSRMQueries struct {
	baseUrl string
	params string
}

/*
 * +----------------------------------------------+
 * | Functions for configure connection to server |
 * +----------------------------------------------+
*/

func (osrm *OSRMQueries) ConfigureOSRMQueries(conn, overview, alternatives, steps string) {
	if (overview != "simplified" && overview != "full" && overview != "false") {
		log.Fatalf("Invalid overview parameter: %s", overview)
	}

	if (alternatives != "true" && alternatives != "false") {
		log.Fatalf("Invalid alternatives parameter: %s", alternatives)
	}

	if (steps != "true" && steps != "false") {
		log.Fatalf("Invalid steps parameter: %s", steps)
	}

	osrm.baseUrl = fmt.Sprintf("http://%s/route/v1/driving/", conn)

	params := url.Values{}
    params.Add("overview", overview)        // полная геометрия маршрута
    params.Add("alternatives", alternatives)   // без альтернативных маршрутов
    params.Add("steps", steps)           // с пошаговой инструкцией
    params.Add("geometries", "geojson")

	osrm.params = params.Encode()
}

/*
 * +-----------------------------------+
 * | Functions for working with routes |
 * +-----------------------------------+
*/

func (osrm *OSRMQueries) BuildRoute(points *dto.CreateRouteDTO) (*dto.RouteDTO, error) {
    // Формируем URL с координатами
	var coordStrings []string
    for _, coord := range points.Stops {
        if len(coord) >= 2 {
            coordStrings = append(coordStrings, fmt.Sprintf("%f,%f", coord[1], coord[0]))
        }
    }
    
    coordinatesStr := strings.Join(coordStrings, ";")
    
    fullURL := fmt.Sprintf("%s%s?%s", osrm.baseUrl, coordinatesStr, osrm.params)
    
    // Выполняем запрос
    resp, err := http.Get(fullURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    // Парсим ответ
    var osrmResponse OSRMResponse
    err = json.NewDecoder(resp.Body).Decode(&osrmResponse)
    if err != nil {
        return nil, err
    }
    
    // Проверяем код ответа
    if osrmResponse.Code != "Ok" {
        return nil, OSRMApiError{fmt.Sprintf("OSRM error: %s", osrmResponse.Code)}
    }

	coordinates := ConvertCoordinatesToLatLng(osrmResponse.Routes[0].Geometry.Coordinates)
    
    return &(dto.RouteDTO{Way: coordinates,
            Distance: osrmResponse.Routes[0].Distance,
            Duration: osrmResponse.Routes[0].Duration}),
            nil
}

/*
 * +---------------------+
 * | Auxiliary functions |
 * +---------------------+
*/

func ConvertCoordinatesToLatLng(coordinates [][]float64) [][]float64 {
	var converted [][]float64
    for _, coord := range coordinates {
        if len(coord) >= 2 {
            converted = append(converted, []float64{coord[1], coord[0]}) // [lat, lng]
        }
    }
    return converted
}