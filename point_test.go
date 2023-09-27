package geo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

// Tests that a call to NewPoint should return a pointer to a Point with the specified values assigned correctly.
func TestNewPoint(t *testing.T) {
	p := NewPoint(40.5, 120.5)

	if p == nil {
		t.Error("Expected to get a pointer to a new point, but got nil instead.")
	}

	if p.Lat != 40.5 {
		t.Errorf("Expected to be able to specify 40.5 as the lat value of a new point, but got %f instead", p.Lat)
	}

	if p.Lon != 120.5 {
		t.Errorf("Expected to be able to specify 120.5 as the lon value of a new point, but got %f instead", p.Lon)
	}
}

// Tests that calling GetLat() after creating a new point returns the expected lat value.
func TestLat(t *testing.T) {
	p := NewPoint(40.5, 120.5)

	lat := p.Lat

	if lat != 40.5 {
		t.Errorf("Expected a call to GetLat() to return the same lat value as was set before, but got %f instead", lat)
	}
}

// Tests that calling GetLng() after creating a new point returns the expected lon value.
func TestLng(t *testing.T) {
	p := NewPoint(40.5, 120.5)

	lng := p.Lon

	if lng != 120.5 {
		t.Errorf("Expected a call to GetLng() to return the same lat value as was set before, but got %f instead", lng)
	}
}

// Seems brittle :\
func TestGreatCircleDistance(t *testing.T) {
	// Test that SEA and SFO are ~ 1091km apart, accurate to 100 meters.
	sea := &Point{Lat: 47.4489, Lon: -122.3094}
	sfo := &Point{Lat: 37.6160933, Lon: -122.3924223}
	sfoToSea := 1093.379199082169

	dist := sea.GreatCircleDistance(sfo)

	if !(dist < (sfoToSea+0.1) && dist > (sfoToSea-0.1)) {
		t.Error("Unnacceptable result.", dist)
	}
}

func TestPointAtDistanceAndBearing(t *testing.T) {
	sea := &Point{Lat: 47.44745785, Lon: -122.308065668024}
	p := sea.PointAtDistanceAndBearing(1090.7, 180)

	// Expected results of transposing point
	// ~1091km at bearing of 180 degrees
	resultLat := 37.638557
	resultLng := -122.308066

	withinLatBounds := p.Lat < resultLat+0.001 && p.Lat > resultLat-0.001
	withinLngBounds := p.Lon < resultLng+0.001 && p.Lon > resultLng-0.001
	if !(withinLatBounds && withinLngBounds) {
		t.Error("Unnacceptable result.", fmt.Sprintf("[%f, %f]", p.Lat, p.Lon))
	}
}

func TestBearingTo(t *testing.T) {
	p1 := &Point{Lat: 40.7486, Lon: -73.9864}
	p2 := &Point{Lat: 0.0, Lon: 0.0}
	bearing := p1.BearingTo(p2)

	// Expected bearing 60 degrees
	resultBearing := 100.610833

	withinBearingBounds := bearing < resultBearing+0.001 && bearing > resultBearing-0.001
	if !withinBearingBounds {
		t.Error("Unnacceptable result.", fmt.Sprintf("%f", bearing))
	}
}

func TestMidpointTo(t *testing.T) {
	p1 := &Point{Lat: 52.205, Lon: 0.119}
	p2 := &Point{Lat: 48.857, Lon: 2.351}

	p := p1.MidpointTo(p2)

	// Expected midpoint 50.5363°N, 001.2746°E
	resultLat := 50.53632
	resultLng := 1.274614

	withinLatBounds := p.Lat < resultLat+0.001 && p.Lat > resultLat-0.001
	withinLngBounds := p.Lon < resultLng+0.001 && p.Lon > resultLng-0.001
	if !(withinLatBounds && withinLngBounds) {
		t.Error("Unnacceptable result.", fmt.Sprintf("[%f, %f]", p.Lat, p.Lon))
	}
}

// Ensures that a point can be marhalled into JSON
func TestMarshalJSON(t *testing.T) {
	p := NewPoint(40.7486, -73.9864)
	res, err := json.Marshal(p)

	if err != nil {
		log.Print(err)
		t.Error("Should not encounter an error when attempting to Marshal a Point to JSON")
	}

	if string(res) != `{"lat":40.7486,"lon":-73.9864}` {
		t.Error("Point should correctly Marshal to JSON")
	}
}

// Ensures that a point can be unmarhalled from JSON
func TestUnmarshalJSON(t *testing.T) {
	data := []byte(`{"lat":40.7486,"lon":-73.9864}`)
	p := &Point{}
	err := p.UnmarshalJSON(data)

	if err != nil {
		t.Errorf("Should not encounter an error when attempting to Unmarshal a Point from JSON")
	}

	if p.Lat != 40.7486 || p.Lon != -73.9864 {
		t.Errorf("Point has mismatched data after Unmarshalling from JSON")
	}
}

// Ensure that a point can be marshalled into slice of binaries
func TestMarshalBinary(t *testing.T) {
	lat, long := 40.7486, -73.9864
	p := NewPoint(lat, long)
	actual, err := p.MarshalBinary()
	if err != nil {
		t.Error("Should not encounter an error when attempting to Marshal a Point to binary", err)
	}

	expected, err := coordinatesToBytes(lat, long)
	if err != nil {
		t.Error("Unable to convert coordinates to bytes slice.", err)
	}

	if !bytes.Equal(actual, expected) {
		t.Errorf("Point should correctly Marshal to Binary.\nExpected %v\nBut got %v", expected, actual)
	}
}

// Ensure that a point can be unmarshalled from a slice of binaries
func TestUnmarshalBinary(t *testing.T) {
	lat, long := 40.7486, -73.9864
	coordinates, err := coordinatesToBytes(lat, long)
	if err != nil {
		t.Error("Unable to convert coordinates to bytes slice.", err)
	}

	actual := &Point{}
	err = actual.UnmarshalBinary(coordinates)
	if err != nil {
		t.Error("Should not encounter an error when attempting to Unmarshal a Point from binary", err)
	}

	expected := NewPoint(lat, long)
	if !assertPointsEqual(actual, expected, 4) {
		t.Errorf("Point should correctly Marshal to Binary.\nExpected %+v\nBut got %+v", expected, actual)
	}
}

func coordinatesToBytes(lat, long float64) ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, lat); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, long); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Asserts true when the latitude and longtitude of p1 and p2 are equal up to a certain number of decimal places.
// Precision is used to define that number of decimal places.
func assertPointsEqual(p1, p2 *Point, precision int) bool {
	roundedLat1, roundedLng1 := int(p1.Lat*float64(precision))/precision, int(p1.Lon*float64(precision))/precision
	roundedLat2, roundedLng2 := int(p2.Lat*float64(precision))/precision, int(p2.Lon*float64(precision))/precision
	return roundedLat1 == roundedLat2 && roundedLng1 == roundedLng2
}
