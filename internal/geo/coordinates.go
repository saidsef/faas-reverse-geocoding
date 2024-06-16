package geo

// Coordinates defines a structure for geographical coordinates with latitude and longitude.
// It is designed to be used in applications that require geographical locations to be represented
// in a structured format. The latitude (Lat) and longitude (Long) are stored as strings to accommodate
// various formats, but they typically represent decimal degrees.
type Coordinates struct {
	Lat  string `json:"lat"`
	Long string `json:"lon"`
}
