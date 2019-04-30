package graph

import (
	"fmt"
	"math"

	svg "github.com/ajstarks/svgo"
	"github.com/miltfra/lisa-rennt/internal"
	"github.com/miltfra/lisa-rennt/internal/terrain"
)

// A Graph contains the paths from each Point
// in the terrain to every other point. It's 0
// when the path doesn't exist
type Graph struct {
	Trrn    *terrain.Terrain
	pt2indx map[*internal.Point]uint16
	indx2pt map[uint16]*internal.Point
	Mtrx    []float64
	N       int
	path    []uint16
	d       *drawer
}

// PrintInformation prints out all necessary information for the user.
func (g *Graph) PrintInformation() {
	dist := float64(0)
	tY := g.indx2pt[g.path[len(g.path)-1]].X * ry
	for i := 1; i < len(g.path); i++ {
		dist += g.Mtrx[int(g.path[i-1])*g.N+int(g.path[i])]
	}
	busT := tY / velB
	lisT := dist / velL
	daySeconds := int(7*3600 + 30*60 + busT - lisT) // 7:30 - time delta
	// We do not need to worry about negative times because lisa wants to get up
	// late and that should not implicate getting up before  the day has started
	hours := daySeconds / 3600
	minutes := (daySeconds % 3600) / 60
	seconds := daySeconds % 60
	fmt.Println("Lisa...")
	// Start Time
	fmt.Printf("...muss %v:%v:%v Uhr losgehen...\n", hours, minutes, seconds)
	daySeconds = int(7*3600 + 30*60 + busT)
	hours = daySeconds / 3600
	minutes = (daySeconds % 3600) / 60
	seconds = daySeconds % 60
	// Target Time
	fmt.Printf("...wird den Bus %v:%v:%v Uhr treffen.\n", hours, minutes, seconds)
	// y-coordinate
	fmt.Printf("...wird den Bus %.2f Meter von der Haltestelle entfernt treffen.\n", tY)
	// duration and distance of travel
	fmt.Printf("...muss %.2f Meter laufen und benötigt %.2f Sekunden.\n", dist, lisT)
	// coordinates and polygons of path points
	fmt.Println("...muss folgende Wegpunkte passieren:\n")
	for _, i := range g.path {
		fmt.Printf("Polygon: %v, Punkt: [x: %.2f, y: %.2f]\n", g.getPolygonID(i), g.indx2pt[i].X, g.indx2pt[i].Y)
	}
	fmt.Printf("Polygon: %v, Punkt: [x: %.2f, y: %.2f]\n", "S", float64(0), tY)
}

func (g *Graph) getPolygonID(pointIndex uint16) string {
	if pointIndex == 0 {
		return "L"
	}
	for i := range g.Trrn.Plygns {
		for _, p := range g.Trrn.Plygns[i].Corners {
			if p.Point == g.indx2pt[pointIndex] {
				return fmt.Sprintf("P%v", i+1)
			}
		}
	}
	return "P0"
}

// New returns a new Graph based on a terrain
func New(t *terrain.Terrain) *Graph {
	pt2indx := make(map[*internal.Point]uint16)
	indx2pt := make(map[uint16]*internal.Point)
	pt2indx[t.Start] = 0
	indx2pt[0] = t.Start
	count := 1
	for _, plygn := range t.Plygns {
		for _, c := range plygn.Corners {
			if !c.IsConcave() {
				indx2pt[uint16(count)] = c.Point
				pt2indx[c.Point] = uint16(count)
				count++
			}
		}
	}
	mtrx := make([]float64, count*count)
	var d float64
	for i := uint16(0); int(i) < count; i++ {
		for j := uint16(0); j < i; j++ {
			p1 := indx2pt[i]
			p2 := indx2pt[j]
			if !HasIntersections(t, internal.NewLineSegment(*p1, *p2)) {
				d = math.Sqrt(internal.GetSqDist(p1, p2))
			} else {
				d = -1
			}
			mtrx[int(i)*count+int(j)] = d
			mtrx[int(j)*count+int(i)] = d
		}
	}
	return &Graph{t, pt2indx, indx2pt, mtrx, int(count), nil, nil}
}

// HasIntersections is true if the given line segment intersects
// at least one polygon of the terrain
func HasIntersections(t *terrain.Terrain, ls *internal.LineSegment) bool {
	for _, plygn := range t.Plygns {
		if plygn.DoesIntersect(ls) {
			return true
		}
	}
	return false
}

// getMax returns the maximum x and y coordinates of the
// graph
func (g *Graph) getMax() (float64, float64) {
	max := internal.NewPoint(0, 0)
	for i := uint16(0); int(i) < g.N; i++ {
		if g.indx2pt[i].X > max.X {
			max.X = g.indx2pt[i].X
		}
		if g.indx2pt[i].Y > max.Y {
			max.Y = g.indx2pt[i].Y
		}
	}
	if g.path != nil {
		final := g.indx2pt[g.path[len(g.path)-1]]
		if final.X*ry+final.Y > max.Y {
			max.Y = final.X*ry + final.Y
		}
	}
	return float64(int(max.X) + 1), float64(int(max.Y) + 1)
}

// getMaxInt returns the maximum x and y coordinates of the
// graph rounded up.
func (g *Graph) getMaxInt() (int, int) {
	maxX, maxY := g.getMax()
	if maxX > float64(int(maxX)) {
		maxX++
	}
	if maxY > float64(int(maxY)) {
		maxY++
	}
	return int(maxX), int(maxY)
}

// DrawAll executes all draw functions on the SVG in the right
// order.
func (g *Graph) DrawAll(svg *svg.SVG) {
	g.d.DrawInit(svg)
	g.d.drawPossible(svg)
	g.d.DrawObstacles(svg)
	g.d.DrawPath(svg)
	g.d.DrawHome(svg)
	g.d.DrawStreet(svg)
}

// DrawGraph executes all draw functions on the SVG in the right
// order.
func (g *Graph) DrawGraph(svg *svg.SVG) {
	g.d.DrawInit(svg)
	g.d.drawPossible(svg)
	g.d.DrawObstacles(svg)
	g.d.DrawHome(svg)
	g.d.DrawStreet(svg)
}

// DrawPath executes all draw functions on the SVG in the right
// order.
func (g *Graph) DrawPath(svg *svg.SVG) {
	g.d.DrawInit(svg)
	g.d.DrawObstacles(svg)
	g.d.DrawPath(svg)
	g.d.DrawHome(svg)
	g.d.DrawStreet(svg)
}

// DrawObstacles executes all draw functions on the SVG in the right
// order.
func (g *Graph) DrawObstacles(svg *svg.SVG) {
	g.d.DrawInit(svg)
	g.d.DrawObstacles(svg)
	g.d.DrawHome(svg)
	g.d.DrawStreet(svg)
}
