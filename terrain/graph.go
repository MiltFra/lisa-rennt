package terrain

import (
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

// A Graph contains the paths from each Point
// in the terrain to every other point. It's 0
// when the path doesn't exist
type Graph struct {
	Trrn    *Terrain
	pt2indx map[Point]int
	indx2pt map[int]Point
	Mtrx    []float64
	N       int
}

// NewGraph returns a new Graph based on a terrain
func NewGraph(t *Terrain) *Graph {
	pt2indx := make(map[Point]int)
	indx2pt := make(map[int]Point)
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
	for i := 0; i < count; i++ {
		for j := 0; j < i; j++ {
			p1 := indx2pt[i]
			p2 := indx2pt[j]
			if !HasIntersections(t, NewLineSegment(p1, p2)) {
				d := math.Sqrt(GetSqDist(p1, p2))
				mtrx[i*count+j] = d
				mtrx[j*count+i] = d
			}
		}
	}
	return &Graph{t, pt2indx, indx2pt, mtrx, count}
}

// GetAccessiblePoints returns a slice of Points which
// should be considered in pathing
func oldGetAccessiblePoints(t *Terrain, p Point) []Point {
	res := make([]Point, 0, 10)
	var ls *LineSegment
	for _, plygn := range t.Plygns {
		first, last, err := plygn.getShadow(p)
		if err != nil {
			res = append(res, plygn.getUnblockedSpikes(p)...)
		} else {
			// TODO Move this to another function
			var found bool
			var accP1, accP2 Point
			mark := first
			for true {
				if !mark.IsConcave() {
					ls = NewLineSegment(p, mark.Pos)
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

// GetAccessiblePoints returns a slice of Points which
// should be considered in pathing
func GetAccessiblePoints(t *Terrain, p Point) []Point {
	res := make([]Point, 0, 10)
	for _, plygn := range t.Plygns {
		res = append(res, plygn.getUnblockedSpikes(p)...)
	}
	return res
}

// HasIntersections is true if the given line segment intersects
// at least one polygon of the terrain
func HasIntersections(t *Terrain, ls *LineSegment) bool {
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
	pathCol := color.RGBA{240, 50, 50, 255}
	for i := 0; i < G.N; i++ {
		for j := i + 1; j < G.N; j++ {
			if G.Mtrx[i*G.N+j] != 0 {
				x0, y0 := scale(G.indx2pt[i].X, G.indx2pt[i].Y)
				x1, y1 := scale(G.indx2pt[j].X, G.indx2pt[j].Y)
				drawLine(x0, y0, x1, y1, img, pathCol)
			}
		}
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

func drawPolygon(p *Polygon, img *image.RGBA64, col color.Color) {
	min, max := p.GetBox()
	for i := min.Y - 1; i <= max.Y+1; i++ {
		nodes := make([]float64, 0, 20)
		//ls := NewLineSegment(NewPoint(min.X-1, i), NewPoint(max.X+1, i))
		for _, c := range p.Corners {
			if (c.Pos.Y <= i) != (c.N2.Pos.Y <= i) {
				nodes = append(nodes, c.Pos.X+(i-c.Pos.Y)/(c.N2.Pos.Y-c.Pos.Y)*(c.N2.Pos.X-c.Pos.X))
			}
		}
		sort(nodes)
		for k := 0; k < len(nodes)-1; k += 2 {
			x0, y0 := scale(nodes[k], i)
			x1, y1 := scale(nodes[k+1], i)
			drawLine(x0, y0, x1, y1, img, col)

		}
	}

}

func sort(arr []float64) {
	var min int
	for i := 0; i < len(arr); i++ {
		min = i
		for j := i; j < len(arr); j++ {
			if arr[min] > arr[j] {
				min = j
			}
		}
		arr[i], arr[min] = arr[min], arr[i]
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
	max := NewPoint(0, 0)
	for i := 0; i < G.N; i++ {
		if G.indx2pt[i].X > max.X {
			max.X = G.indx2pt[i].X
		}
		if G.indx2pt[i].Y > max.Y {
			max.Y = G.indx2pt[i].Y
		}
	}
	return float64(int(max.X) + 1), float64(int(max.Y) + 1)
}
