package db

import (
	"GruzowikiRoutesGenerator/utils/dto"
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	"github.com/jackc/pgx/v5"
)

type PGRoutingQueries struct {
	Connection *pgx.Conn
}

/*
 * +---------------------------------------+
 * | Functions for working with connection |
 * +---------------------------------------+
*/

func (pgr *PGRoutingQueries) EstablishConnection(connString string) {
	var err error

	pgr.Connection, err = pgx.Connect(context.Background(), connString);
	if err != nil {
		log.Fatal(err)
	}
}

func (pgr *PGRoutingQueries) FuckingDestroyConnection() {
	pgr.Connection.Close(context.Background())
}

/*
 * +-----------------------------------+
 * | Functions for working with routes |
 * +-----------------------------------+
*/

func (pgr *PGRoutingQueries) BuildRoute(points *dto.CreateRouteDTO) (dto.RouteDTO, error) {
	var vertices *list.List = list.New();

	var vertexId int64
	for _, point := range points.Stops {
		query := fmt.Sprintf(
			`SELECT id
			FROM routing_roads_vertices_pgr
			WHERE ST_DWithin(
				geom,
				ST_SetSRID(ST_Point(%f, %f), 4326),
				100
			)
			ORDER BY geom <-> ST_SetSRID(ST_Point(%f, %f), 4326)
			LIMIT 1`, point[1], point[0], point[1], point[0])

		err := pgr.Connection.QueryRow(context.Background(), query).Scan(&vertexId);
		if err != nil {
			log.Error("Map geo-point to graph vertex failed: %v\n", err)
			return dto.RouteDTO{}, MapPointsError{"Map geo-point to graph vertex failed"}
		}

		vertices.PushBack(vertexId)
	}

	if vertices.Len() != len(points.Stops) {
		log.Error("Can not create route of one vertex")
		return dto.RouteDTO{}, MapPointsError{"Can not create route of one vertex"}
	}

	verticesString := "ARRAY" + ListToString(*vertices)

	query := fmt.Sprintf(
		`WITH dijkstra AS (
			SELECT *
			FROM pgr_dijkstraVia(
				'SELECT id, source, target, cost, reverse_cost FROM routing_roads',
				%s,
				directed := false
			)
		),
		path_edges AS (
			SELECT 
				d.path_id AS PathId,
				d.path_seq AS Sequence,
				d.edge AS EdgeId,
				d.node AS NodeId,
				r.geom AS edge_geom
			FROM dijkstra d
			LEFT JOIN routing_roads r ON d.edge = r.id
			WHERE d.edge > 0
		)
		SELECT 
			PathId,
			ST_AsGeoJSON(ST_FlipCoordinates(ST_LineMerge(ST_Collect(edge_geom)))) AS RouteGeometry,
			ST_Length(ST_LineMerge(ST_Collect(edge_geom))) AS RouteLength
		FROM path_edges
		GROUP BY PathId
		ORDER BY PathId`, verticesString)
	
	var pathId int64
	var routeGeomString string
	var routeLength float64
	err := pgr.Connection.QueryRow(context.Background(), query).Scan(&pathId, &routeGeomString, &routeLength);
	if err != nil {
		log.Error("Build road failed: %v\n", err)
		return dto.RouteDTO{Way: [][]float64{}, Graph: []int32{}}, BuildRouteError{"Build road failed"}
	}

	var routeGeom PgrGeoJson
	err = json.Unmarshal([]byte(routeGeomString), &routeGeom)
	if err != nil {
		log.Error("Failed to parse JSON: %v\n", err)
		return dto.RouteDTO{Way: [][]float64{}, Graph: []int32{}}, BuildRouteError{"JSON parse failed"}
	}

	return dto.RouteDTO{Way: routeGeom.Coordinates, Graph: []int32{}}, nil
}

/*
 * +---------------------+
 * | Auxiliary functions |
 * +---------------------+
*/

func ListToString(list list.List) string {
	var builder strings.Builder = strings.Builder{}
	builder.WriteString("[")

	item := list.Front()
	for ; ; {
		switch val := item.Value.(type) {
		case int64:
			_, _ = builder.WriteString(strconv.FormatInt(val, 10))
		default:
			log.Warn("Can not convert list item to string")
			return ""
		}
		

		item = item.Next()
		if item == nil {
			builder.WriteString("]")
			break
		} else {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}