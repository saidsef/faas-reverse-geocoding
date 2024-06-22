package geo

import (
	"encoding/json"
	"strconv"
)

// Coordinates defines a structure for geographical coordinates with latitude and longitude.
// It is designed to be used in applications that require geographical locations to be represented
// in a structured format. The latitude (Lat) and longitude (Long) are stored as strings to accommodate
// various formats, but they typically represent decimal degrees.
type Coordinates struct {
	Lat  float32 `json:"lat"`
	Long float32 `json:"lon"`
}

func (c *Coordinates) UnmarshalJSON(data []byte) error {
	var aux struct {
		Lat  string `json:"lat"`
		Long string `json:"lon"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	lat, err := strconv.ParseFloat(aux.Lat, 32)
	if err != nil {
		return err
	}

	long, err := strconv.ParseFloat(aux.Long, 32)
	if err != nil {
		return err
	}

	c.Lat = float32(lat)
	c.Long = float32(long)
	return nil
}
