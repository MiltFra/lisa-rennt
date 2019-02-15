package terrain

import (
	"errors"
	"math"
)

// A Line represents a line in the 2d-Plane:
// A * x + B * y = C
type Line struct {
	A float64
	B float64
	C float64
}

// A LineSegment is a part of a line in certain
// boundaries
type LineSegment struct {
	L *Line
	A Point
	B Point
}

// ErrBadLineOrientation occurs, when it's
// impossible to convert between x and y of a
// line because m is either infinitly large or 0
var ErrBadLineOrientation = errors.New(
	"Line orientation does not allow conversion between X and Y")

// NewLineSegment returns a new line segment
// that connects two given points
func NewLineSegment(A, B Point) *LineSegment {
	// See https://de.wikipedia.org/wiki/Koordinatenform
	L := &Line{A.Y - B.Y, B.X - A.X, B.X*A.Y - A.X*B.Y}
	return &LineSegment{L, A, B}
}

// GetY returns the corresponding y to a given x
// of a line; if possible
func (L *Line) GetY(X float64) (float64, error) {
	if L.B == 0 {
		return 0, ErrBadLineOrientation
	}
	return (L.C - L.A*X) / L.B, nil
}

// GetX returns the corresponding x to a given y
// of a line; if possible
func (L *Line) GetX(Y float64) (float64, error) {
	if L.A == 0 {
		return 0, ErrBadLineOrientation
	}
	return L.C/L.A - (L.B/L.A)*Y, nil
}

// GetMinMax returns the boundaries of the rectangle
// the endpoints of the line shape
func (LS *LineSegment) GetMinMax() (Point, Point) {
	return sortPointCoordinates(LS.A, LS.B)
}

// RelationOf is true, when the given point is
// below the line, hence the line is above the
// point. For the sake of simplicity above means left
// of a vertical line (thus you can check wether
// points are on the same side of the line)
func (L *Line) RelationOf(P Point) int {
	if L.B == 0 {
		if P.X > L.C/L.A {
			return 1 // Above
		} else if P.X < L.C/L.A {
			return -1 // Below
		} else {
			return 0 // On
		}
	}
	y, _ := L.GetY(P.X)
	if y > P.Y {
		return 1 // Left
	} else if y < P.Y {
		return -1 // Right
	} else {
		return 0 // On
	}
}

// intersection returns true when the path A
// crosses the path B
// NOTE: having one or more common points
// does not impy an intersection!
func intersection(L1, L2 *LineSegment) int {
	rel1 := L1.L.RelationOf(L2.A) * L1.L.RelationOf(L2.B)
	rel2 := L2.L.RelationOf(L1.A) * L2.L.RelationOf(L1.B)
	if rel1 > 0 {
		return 1
	}
	if rel2 > 0 {
		return 1
	}
	if rel1 == 0 {
		return 0
	}
	if rel2 == 0 {
		return 0
	}
	return -1
}

// sortPointCoordinates returns the given points in such a way
// that no coordinate of P is bigger than its counterpart in Q
func sortPointCoordinates(P, Q Point) (Point, Point) {
	if P.X > Q.X {
		P.X, Q.X = Q.X, P.X
	}
	if P.Y > Q.Y {
		P.Y, Q.Y = Q.Y, P.Y
	}
	return P, Q
}

// IsConcave returns true if the angle at this
// corner is greater than 180°
// See https://bit.ly/2xm22oj
func (crnr *Corner) IsConcave() bool {
	a := crnr.N1.Pos
	b := crnr.Pos
	c := crnr.N2.Pos
	return (b.X-a.X)*(c.Y-b.Y)-(b.Y-a.Y)*(c.X-b.X) < 0
}

// passes returns true if a given line segment intersects
// the borders of a given box or has a point inside the box
func (LS *LineSegment) passes(min, max Point) bool {
	lsMin, lsMax := LS.GetMinMax()
	if min.X > lsMax.X || lsMin.X > max.X ||
		min.Y > lsMax.Y || lsMin.Y > max.Y {
		return false
	}
	//if LS.A.IsInside(min, max) || LS.B.IsInside(min, max) {
	//	return true
	//}
	l1 := NewLineSegment(min, Point{min.X, max.Y})
	if intersection(LS, l1) > 0 {
		return true
	}
	l2 := NewLineSegment(min, Point{max.X, min.Y})
	if intersection(LS, l2) > 0 {
		return true
	}
	l3 := NewLineSegment(Point{min.X, max.Y}, max)
	if intersection(LS, l3) > 0 {
		return true
	}
	l4 := NewLineSegment(Point{max.X, min.Y}, max)
	if intersection(LS, l4) > 0 {
		return true
	}
	return false
}

// ScaleTo adjusts the Point B of the line segment in
// such a way that the length of the line becomes s and
// the direction is kept
func (LS *LineSegment) ScaleTo(s float64) *LineSegment {
	r := s / math.Sqrt(math.Pow(LS.A.X-LS.B.X, 2)+math.Pow(LS.A.Y-LS.B.Y, 2))
	LS.B = NewPoint(LS.A.X+(LS.B.X-LS.A.X)*r, LS.A.Y+(LS.B.Y-LS.A.Y)*r)
	return LS
}

func (LS *LineSegment) inCounterClockwiseAngle(p0, p1, p2 Point) bool {
	angle1, err := getAngle(p0, p1)
	if err != nil {
		return false
	}
	angle2, err := getAngle(p0, p2)
	if err != nil {
		return false
	}
	angleA, errA := getAngle(p0, LS.A)
	angleB, errB := getAngle(p0, LS.B)
	condA := angleA > angle1 && angleA < angle2
	condB := angleB > angle1 && angleB < angle2
	if errA == nil && errB == nil {
		return condA != condB
	}
	if errA == nil {
		return condA
	} else if errB == nil {
		return condB
	}
	return false
}

func (LS *LineSegment) changesSide(p0, p1, p2 Point) bool {
	angleA, err := getAngle(p0, LS.A)
	if err != nil {
		return false
	}
	angleB, err := getAngle(p0, LS.B)
	if err != nil {
		return false
	}
	angle1, err := getAngle(p0, p1)
	if err != nil {
		return false
	}
	angle2, err := getAngle(p0, p2)
	if err != nil {
		return false
	}
	// Swapping angles to reduce cases
	if angleA > angleB {
		angleA, angleB = angleB, angleA
	}
	enclosed1 := angleA < angle1 && angleB > angle1
	enclosed2 := angleA < angle2 && angleB > angle2
	return enclosed1 != enclosed2
}
