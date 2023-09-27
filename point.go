package geo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
)

// Point Represents a Physical Point in geographic notation [lat, lon].
type Point struct {
	Lat float64
	Lon float64
}

const (
	// EarthRadius According to Wikipedia, the Earth's radius is about 6,371km
	EarthRadius = 6371
)

// NewPoint returns a new Point populated by the passed in latitude (lat) and longitude (lon) values.
func NewPoint(lat, lon float64) *Point {
	return &Point{Lat: lat, Lon: lon}
}

// PointAtDistanceAndBearing returns a Point populated with the lat and lon coordinates
// by transposing the origin point the passed in distance (in kilometers)
// by the passed in compass bearing (in degrees).
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) PointAtDistanceAndBearing(dist, bearing float64) *Point {

	dr := dist / EarthRadius

	bearing = bearing * (math.Pi / 180.0)

	lat1 := p.Lat * (math.Pi / 180.0)
	lng1 := p.Lon * (math.Pi / 180.0)

	lat2Part1 := math.Sin(lat1) * math.Cos(dr)
	lat2Part2 := math.Cos(lat1) * math.Sin(dr) * math.Cos(bearing)

	lat2 := math.Asin(lat2Part1 + lat2Part2)

	lon2Part1 := math.Sin(bearing) * math.Sin(dr) * math.Cos(lat1)
	lon2Part2 := math.Cos(dr) - (math.Sin(lat1) * math.Sin(lat2))

	lng2 := lng1 + math.Atan2(lon2Part1, lon2Part2)
	lng2 = math.Mod(lng2+3*math.Pi, 2*math.Pi) - math.Pi

	lat2 = lat2 * (180.0 / math.Pi)
	lng2 = lng2 * (180.0 / math.Pi)

	return &Point{Lat: lat2, Lon: lng2}
}

// GreatCircleDistance Calculates the Haversine distance between two points in kilometers.
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) GreatCircleDistance(p2 *Point) float64 {
	dLat := (p2.Lat - p.Lat) * (math.Pi / 180.0)
	dLon := (p2.Lon - p.Lon) * (math.Pi / 180.0)

	lat1 := p.Lat * (math.Pi / 180.0)
	lat2 := p2.Lat * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadius * c
}

// BearingTo Calculates the initial bearing (sometimes referred to as forward azimuth)
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) BearingTo(p2 *Point) float64 {

	dLon := (p2.Lon - p.Lon) * math.Pi / 180.0

	lat1 := p.Lat * math.Pi / 180.0
	lat2 := p2.Lat * math.Pi / 180.0

	y := math.Sin(dLon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) -
		math.Sin(lat1)*math.Cos(lat2)*math.Cos(dLon)
	brng := math.Atan2(y, x) * 180.0 / math.Pi

	return brng
}

// MidpointTo Calculates the midpoint between 'this' point and the supplied point.
// Original implementation from http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) MidpointTo(p2 *Point) *Point {
	lat1 := p.Lat * math.Pi / 180.0
	lat2 := p2.Lat * math.Pi / 180.0

	lon1 := p.Lon * math.Pi / 180.0
	dLon := (p2.Lon - p.Lon) * math.Pi / 180.0

	bx := math.Cos(lat2) * math.Cos(dLon)
	by := math.Cos(lat2) * math.Sin(dLon)

	lat3Rad := math.Atan2(
		math.Sin(lat1)+math.Sin(lat2),
		math.Sqrt(math.Pow(math.Cos(lat1)+bx, 2)+math.Pow(by, 2)),
	)
	lon3Rad := lon1 + math.Atan2(by, math.Cos(lat1)+bx)

	lat3 := lat3Rad * 180.0 / math.Pi
	lon3 := lon3Rad * 180.0 / math.Pi

	return NewPoint(lat3, lon3)
}

// MarshalBinary renders the current point to a byte slice.
// Implements the encoding.BinaryMarshaler Interface.
func (p *Point) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, p.Lat)
	if err != nil {
		return nil, fmt.Errorf("unable to encode lat %v: %v", p.Lat, err)
	}
	err = binary.Write(&buf, binary.LittleEndian, p.Lon)
	if err != nil {
		return nil, fmt.Errorf("unable to encode lon %v: %v", p.Lon, err)
	}

	return buf.Bytes(), nil
}

func (p *Point) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	var lat float64
	err := binary.Read(buf, binary.LittleEndian, &lat)
	if err != nil {
		return fmt.Errorf("binary.Read failed: %v", err)
	}

	var lng float64
	err = binary.Read(buf, binary.LittleEndian, &lng)
	if err != nil {
		return fmt.Errorf("binary.Read failed: %v", err)
	}

	p.Lat = lat
	p.Lon = lng
	return nil
}

// MarshalJSON renders the current Point to valid JSON.
// Implements the json.Marshaller Interface.
func (p *Point) MarshalJSON() ([]byte, error) {
	res := fmt.Sprintf(`{"lat":%v, "lon":%v}`, p.Lat, p.Lon)
	return []byte(res), nil
}

// UnmarshalJSON decodes the current Point from a JSON body.
// Throws an error if the body of the point cannot be interpreted by the JSON body
func (p *Point) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	var values map[string]float64
	err := dec.Decode(&values)

	if err != nil {
		return err
	}

	*p = *NewPoint(values["lat"], values["lon"])

	return nil
}
