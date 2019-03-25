package graph

import (
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"sort"

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
	return &Graph{t, pt2indx, indx2pt, mtrx, count, nil}
}

// GetAccessiblePoints returns a slice of Points which
// should be considered in pathing
func oldGetAccessiblePoints(t *terrain.Terrain, p *internal.Point) []*internal.Point {
	res := make([]*internal.Point, 0, 10)
	var ls *internal.LineSegment
	for _, plygn := range t.Plygns {
		first, last, err := plygn.GetShadow(p)
		if err != nil {
			res = append(res, plygn.GetUnblockedSpikes(p)...)
		} else {
			// TODO Move this to another function
			var found bool
			var accP1, accP2 *internal.Point
			mark := first
			for true {
				if !mark.IsConcave() {
					ls = internal.NewLineSegment(*p, *mark.Pos)
					if !HasIntersections(t, ls) {
						if !found {
							accP1 = mark.Pos
							found = true
						}
						accP2 = mark.Pos
					}
				}
				if mark != last {
					mark = mark.N1
				} else {
					break
				}
			}
			if found {
				res = append(res, accP1)
				if accP1 != accP2 {
					res = append(res, accP2)
				}
			}
		}
	}
	return res
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

var imgX float64 = 1920
var imgY float64 = 1000
var maxX, maxY float64

// Draw draws the graph to an image
func (G *Graph) Draw(path string) {
	maxX, maxY = G.getMax()
	img := image.NewRGBA64(image.Rect(0, 0, int(imgX), int(imgY)))
	// Colors
	//fillCol := color.RGBA{100, 100, 100, 255}
	// Draw Polygons
	plygnCol := color.RGBA{255, 255, 255, 255}
	for _, plygn := range G.Trrn.Plygns {
		//drawPolygon(plygn, img, fillCol)
		for _, c := range plygn.Corners {
			x0, y0 := scale(c.Pos.X, c.Pos.Y)
			x1, y1 := scale(c.N2.Pos.X, c.N2.Pos.Y)
			drawLine(x0, y0, x1, y1, img, plygnCol)
		}
	}
	if G.path == nil {
		pathCol := color.RGBA{240, 50, 50, 255}
		for i := 0; i < G.N; i++ {
			for j := i + 1; j < G.N; j++ {
				if G.Mtrx[i*G.N+j] != -1 {
					x0, y0 := scale(G.indx2pt[i].X, G.indx2pt[i].Y)
					x1, y1 := scale(G.indx2pt[j].X, G.indx2pt[j].Y)
					drawLine(x0, y0, x1, y1, img, pathCol)
				}
			}
		}
	} else {
		shortestCol := color.RGBA{50, 50, 240, 255}
		for i := 1; i < len(G.path); i++ {
			p1 := G.indx2pt[G.path[i-1]]
			p2 := G.indx2pt[G.path[i]]
			x0, y0 := scale(p1.X, p1.Y)
			x1, y1 := scale(p2.X, p2.Y)
			drawLine(x0, y0, x1, y1, img, shortestCol)
		}
		p0 := G.indx2pt[G.path[len(G.path)-1]]
		p1 := internal.NewPoint(0, p0.Y+p0.X*ry)
		x0, y0 := scale(p0.X, p0.Y)
		x1, y1 := scale(p1.X, p1.Y)
		drawLine(x0, y0, x1, y1, img, shortestCol)

	}
	toimg, _ := os.Create(path)
	defer toimg.Close()
	jpeg.Encode(toimg, img, nil)
}

func scale(x, y float64) (int, int) {
	factor := math.Min(imgX/maxX, imgY/maxY)
	return int(x * factor), int(y * factor)
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func drawPolygon(p *internal.Polygon, img *image.RGBA64, col color.Color) {
	min, max := p.GetBox()
	for i := min.Y - 1; i <= max.Y+1; i++ {
		nodes := make([]float64, 0, 20)
		//ls := NewLineSegment(NewPoint(min.X-1, i), NewPoint(max.X+1, i))
		for _, c := range p.Corners {
			if (c.Pos.Y <= i) != (c.N2.Pos.Y <= i) {
				nodes = append(nodes, c.Pos.X+(i-c.Pos.Y)/(c.N2.Pos.Y-c.Pos.Y)*(c.N2.Pos.X-c.Pos.X))
			}
		}
		sort.Float64s(nodes)
		for k := 0; k < len(nodes)-1; k += 2 {
			x0, y0 := scale(nodes[k], i)
			x1, y1 := scale(nodes[k+1], i)
			drawLine(x0, y0, x1, y1, img, col)
		}
	}
}

// drawLine draws a line between two points on an image
func drawLine(x0, y0, x1, y1 int, img *image.RGBA64, col color.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy
	for {
		img.Set(x0, y0, col)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
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
