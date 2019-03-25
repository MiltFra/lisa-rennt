package internal

import "math"

// A Point is a pair of coordinates in the
// 2D-plane
type Point struct {
	X float64 // float even though the coordinates will be in uint range
	Y float64 // so less conversions are needed
}

// NewPoint returns a new point with coordinates
// x and y
func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

// GetAngle returns the absolute angle of a point as seen from
// another point in a cartesian coordinate system.
//
// Values are in [0; 360[ in degree.
//
// 0 means the point has the same x and is on the left of the
// pov. Any other value is counter clockwise to this position.
func GetAngle(pov, point *Point) (float64, error) {
	if point.X == pov.X {
		if point.Y == pov.Y {
			return 0, ErrUndefinedAngle
		} else if point.Y > pov.Y {
			return math.Pi / 2, nil
		}
		return math.Pi * 3 / 2, nil
	}
	if point.Y == pov.Y {
		if point.X > pov.X {
			return 0, nil
		}
		return math.Pi, nil
	}
	angle := math.Atan((point.Y - pov.Y) / (point.X - pov.X))
	// If the point is on the left of the point we need
	// to add 180 degrees
	if point.X < pov.X {
		angle += math.Pi
	}
	// Since our points are now in the interval [-90; 270)
	// we need to fit it to [0; 360)
	return ToPositiveAngle(angle), nil
}

// IsInside is true if the point is inside the given
// rectangle
func (P *Point) IsInside(min, max Point) bool {
	return min.X > P.X && P.X < max.X && min.Y > P.Y && P.Y < max.Y
}

// GetSqDist returns the square of the distance of two
// given points; this allows for distance comparison
// without the sqrt function
func GetSqDist(p, q *Point) float64 {
	return math.Pow(p.X-q.X, 2) + math.Pow(p.Y-q.Y, 2)
}

// ToPositiveAngle moves any angle into the interval
// [0, math.Pi * 2)
func ToPositiveAngle(a float64) float64 {
	for a > math.Pi*2 {
		a -= math.Pi * 2
	}
	return a
}
