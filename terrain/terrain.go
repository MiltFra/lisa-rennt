package terrain

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/miltfra/tools"
)

// Terrain describes a certain instance of the
// Lisa Rennt problem
type Terrain struct {
	velBus float64
	velLis float64
	PolCnt int
	Plygns []*Polygon
	Start  Point
	TanA   float64
}

// A Point is a pair of coordinates in the
// 2D-plane
type Point struct {
	X float64 // float even though the coordinates will be in uint range
	Y float64 // so less conversions are needed
}

// NewPoint returns a new point with coordinates
// x and y
func NewPoint(x, y float64) Point {
	return Point{x, y}
}

// New returns a Terrain instance read from a file
func New(path string) *Terrain {
	// Open a stream to the file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("[ERR] Error while reading file:\n", err)
	}
	// Close the stream at the end of this function
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// Read to the next \n
	scanner.Scan()
	polCnt := tools.Stoi(scanner.Text())
	plygns := make([]*Polygon, polCnt)
	for i := 0; i < polCnt; i++ {
		scanner.Scan()
		plygns[i] = stoPlygn(scanner.Text())
	}
	scanner.Scan()
	h := strings.Split(scanner.Text(), " ")
	home := NewPoint(tools.Stof64(h[0]), tools.Stof64(h[1]))
	velB := 30 / 3.6
	velL := 15 / 3.6
	return &Terrain{
		velB, velL, polCnt, plygns, home,
		math.Sqrt((math.Pow(velB, 2) / math.Pow(velL, 2)) - 1)}
}

// IsInside is true if the point is inside the given
// rectangle
func (P *Point) IsInside(min, max Point) bool {
	return min.X > P.X && P.X < max.X && min.Y > P.Y && P.Y < max.Y
}

// ErrUndefinedAngle is returned when the angle between two points
// p and q is undefined. e.g. when p == q
var ErrUndefinedAngle = errors.New("Can't compute between given points")

func getAngle(pov, point Point) (float64, error) {
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

// GetSqDist returns the square of the distance of two
// given points; this allows for distance comparison
// without the sqrt function
func GetSqDist(p, q Point) float64 {
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
