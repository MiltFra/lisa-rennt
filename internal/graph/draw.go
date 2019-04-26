package graph

import (
	svg "github.com/ajstarks/svgo"
	"github.com/miltfra/lisa-rennt/internal"
)

var (
	formatPossbile     = "stroke:#55aadd;stroke-width:3"
	formatObstacle     = "fill:#6b6b6b;stroke:black"
	formatBg           = "fill:#CCFFCC"
	formatShortest     = "stroke:#000080;stroke-width:3;fill:none"
	formatHome         = "stroke:black;fill:red"
	formatStreet       = "stroke:black;fill:#B3B3B3;stroke-width:4"
	formatStreetMiddle = "stroke:black;stroke-dasharray:25,25;stroke-width:3"
	radiusHome         = 10
	widthStreet        = 80
)

type drawer struct {
	maxX, maxY int
	ogX, ogY   int
	invX, invY bool
	g          *Graph
}

func newDrawer(g *Graph) *drawer {
	maxX, maxY := g.getMaxInt()
	ogX, ogY := widthStreet, maxY
	maxX += widthStreet
	invX, invY := false, true
	return &drawer{
		maxX, maxY,
		ogX, ogY,
		invX, invY,
		g,
	}
}

// DrawPossible draws all possible paths onto the
// given SVG. Nothing else is added.
func (d *drawer) drawPossible(canv *svg.SVG) {
	// For every starting node
	for i := 0; i < d.g.N; i++ {
		// For every target node
		for j := i + 1; j < d.g.N; j++ {
			// If this edge exists
			if d.g.Mtrx[i*d.g.N+j] != -1 {
				// Draw the line
				p0 := d.g.indx2pt[i]
				p1 := d.g.indx2pt[j]
				x0, y0 := d.Translate(p0.X, p0.Y)
				x1, y1 := d.Translate(p1.X, p1.Y)
				canv.Line(x0, y0, x1, y1, formatPossbile)
			}
		}
	}
}

// DrawInit creates the background and the polygons on
// the given SVG. Nothin else is added.
func (d *drawer) DrawInit(canv *svg.SVG) {
	canv.Start(int(d.maxX), int(d.maxY))
	canv.Rect(0, 0, int(d.maxX), int(d.maxY), formatBg)

}

// DrawObstacles draws the polygons on the given SVG.
func (d *drawer) DrawObstacles(canv *svg.SVG) {
	for _, plygn := range d.g.Trrn.Plygns {
		x := make([]int, plygn.N)
		y := make([]int, plygn.N)
		for i, c := range plygn.Corners {
			x[i], y[i] = d.Translate(c.X, c.Y)
		}
		canv.Polygon(x, y, formatObstacle)
	}
}

// DrawStreet draws the street on the given SVG.
func (d *drawer) DrawStreet(canv *svg.SVG) {
	x0, y0 := d.Translate(0, 0)
	// Moving the street left to have the entire border on
	// the canvas.
	canv.Rect(2, -10, x0-2, y0+20, formatStreet)
	canv.Line(widthStreet/2, 0, widthStreet/2, y0, formatStreetMiddle)
}

// DrawHome draws the red circle indicating the home
// of lisa at it's coordinates. Nothing else is added.
func (d *drawer) DrawHome(canv *svg.SVG) {
	home := d.g.indx2pt[0]
	x0, y0 := d.Translate(home.X, home.Y)
	canv.Circle(x0, y0, radiusHome, formatHome)
}

// DrawPath draws the path of the graph onto the
// given SVG. Nothing else is added.
func (d *drawer) DrawPath(canv *svg.SVG) {
	x := make([]int, len(d.g.path)+1)
	y := make([]int, len(d.g.path)+1)
	var p *internal.Point
	for i := 0; i < len(d.g.path); i++ {
		p = d.g.indx2pt[d.g.path[i]]
		x[i], y[i] = d.Translate(p.X, p.Y)
	}
	x[len(d.g.path)], y[len(d.g.path)] = d.Translate(0, p.Y+ry*p.X)
	for i := 1; i < len(d.g.path); i++ {
		canv.Polyline(x, y, formatShortest)
	}
}

// Translate returns the coordinates of a point in our
// terrain in the vector graphic.
func (d *drawer) Translate(inX, inY float64) (x, y int) {
	x = int(inX)
	if d.invX {
		x *= -1
	}
	x += d.ogX
	y = int(inY)
	if d.invY {
		y *= -1
	}
	y += d.ogY
	return
}
