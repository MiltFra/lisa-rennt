package graph

import (
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
	pt2indx map[*internal.Point]int
	indx2pt map[int]*internal.Point
	Mtrx    []float64
	N       int
	path    []int
	d       *drawer
}

// New returns a new Graph based on a terrain
func New(t *terrain.Terrain) *Graph {
	pt2indx := make(map[*internal.Point]int)
	indx2pt := make(map[int]*internal.Point)
	pt2indx[t.Start] = 0
	indx2pt[0] = t.Start
	count := 1
	for _, plygn := range t.Plygns {
		for _, c := range plygn.Corners {
			if !c.IsConcave() {
				indx2pt[count] = c.Pos
				pt2indx[c.Pos] = count
				count++
			}
		}
	}
	mtrx := make([]float64, count*count)
	var d float64
	for i := 0; i < count; i++ {
		for j := 0; j < i; j++ {
			p1 := indx2pt[i]
			p2 := indx2pt[j]
			if !HasIntersections(t, internal.NewLineSegment(*p1, *p2)) {
				d = math.Sqrt(internal.GetSqDist(p1, p2))
			} else {
				d = -1
			}
			mtrx[i*count+j] = d
			mtrx[j*count+i] = d
		}
	}
	return &Graph{t, pt2indx, indx2pt, mtrx, count, nil, nil}
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
func (G *Graph) getMax() (float64, float64) {
	max := internal.NewPoint(0, 0)
	for i := 0; i < G.N; i++ {
		if G.indx2pt[i].X > max.X {
			max.X = G.indx2pt[i].X
		}
		if G.indx2pt[i].Y > max.Y {
			max.Y = G.indx2pt[i].Y
		}
	}
	if G.path != nil {
		final := G.indx2pt[G.path[len(G.path)-1]]
		if final.X*ry+final.Y > max.Y {
			max.Y = final.X*ry + final.Y
		}
	}
	return float64(int(max.X) + 1), float64(int(max.Y) + 1)
}

// getMaxInt returns the maximum x and y coordinates of the
// graph rounded up.
func (G *Graph) getMaxInt() (int, int) {
	maxX, maxY := G.getMax()
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
