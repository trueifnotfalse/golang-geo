// Also added other functions and some tests related to geo based polygons.

package geo

// A Polygon is carved out of a 2D plane by a set of (possibly disjoint) contours.
// It can thus contain holes, and can be self-intersecting.
type Polygon struct {
	points []Point
}

// NewPolygon Creates and returns a new pointer to a Polygon
// composed of the passed in points.  Points are
// considered to be in order such that the last point
// forms an edge with the first point.
func NewPolygon(points []Point) *Polygon {
	return &Polygon{points: points}
}

// Points returns the points of the current Polygon.
func (p *Polygon) Points() []Point {
	return p.points
}

// IsClosed returns whether or not the polygon is closed.
// TODO:  This can obviously be improved, but for now,
//        this should be sufficient for detecting if points
//        are contained using the raycast algorithm.
func (p *Polygon) IsClosed() bool {
	if len(p.points) < 3 {
		return false
	}

	return true
}

// Contains returns whether or not the current Polygon contains the passed in Point.
func (p *Polygon) Contains(point *Point) bool {
	if !p.IsClosed() {
		return false
	}

	contains := PNPoly(point, &p.points[len(p.points)-1], &p.points[0])
	for i := 1; i < len(p.points); i++ {
		if PNPoly(point, &p.points[i-1], &p.points[i]) {
			contains = !contains
		}
	}

	return contains
}

func PNPoly(p, a, b *Point) bool {
	return (a.Lon > p.Lon) != (b.Lon > p.Lon) &&
		p.Lat < (b.Lat-a.Lat)*(p.Lon-a.Lon)/(b.Lon-a.Lon)+a.Lat
}
