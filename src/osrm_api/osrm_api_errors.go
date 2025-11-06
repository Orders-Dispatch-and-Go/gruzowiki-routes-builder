package osrmapi

type OSRMApiError struct {
	Text string
}

// Error implements error.
func (m OSRMApiError) Error() string {
	return m.Text
}