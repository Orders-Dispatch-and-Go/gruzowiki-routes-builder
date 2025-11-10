package osrmapi

type OSRMResponse struct {
    Code   string    `json:"code"`
    Routes []Route   `json:"routes"`
    Waypoints []Waypoint `json:"waypoints"`
}

type Route struct {
    Distance float64   `json:"distance"` // в метрах
    Duration float64   `json:"duration"` // в секундах
    Geometry GeoJson    `json:"geometry"` // в формате geojson
    Legs     []Leg     `json:"legs"`
    Weight   float64   `json:"weight"`
    WeightName string  `json:"weight_name"`
}

type Leg struct {
    Distance float64   `json:"distance"`
    Duration float64   `json:"duration"`
    Steps    []Step    `json:"steps"`
    Summary  string    `json:"summary"`
}

type Step struct {
    Distance    float64   `json:"distance"`
    Duration    float64   `json:"duration"`
    Geometry    string    `json:"geometry"`
    Name        string    `json:"name"`
    Instruction string    `json:"instruction"`
    Maneuver    Maneuver  `json:"maneuver"`
}

type Maneuver struct {
    Location []float64 `json:"location"`
    BearingAfter int   `json:"bearing_after"`
    Type      string   `json:"type"`
    Modifier  string   `json:"modifier,omitempty"`
}

type Waypoint struct {
    Name      string    `json:"name"`
    Location  []float64 `json:"location"`
    Distance  float64   `json:"distance"`
}

type GeoJson struct {
	Type string `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}