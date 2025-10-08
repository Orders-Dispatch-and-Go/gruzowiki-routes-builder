package db

type MapPointsError struct {
	Text string
}

type BuildRouteError struct {
	Text string
}

// Error implements error.
func (m MapPointsError) Error() string {
	return m.Text
}

func (m BuildRouteError) Error() string {
	return m.Text
}

