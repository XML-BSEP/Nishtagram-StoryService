package domain

type Location struct {
	Location string
	Latitude float64
	Longitude float64
}

func NewLocation(name string, lat float64, long float64) Location {
	return Location{Location: name, Longitude: long, Latitude: lat}
}