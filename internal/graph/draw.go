package graph

import (
	svg "github.com/ajstarks/svgo"
	"github.com/miltfra/lisa-rennt/internal"
)

func (g *Graph) DrawPossible(canv *svg.SVG) {
	_, maxY := g.getMax()
	for i := 0; i < g.N; i++ {
		for j := i + 1; j < g.N; j++ {
			if g.Mtrx[i*g.N+j] != -1 {
				p0 := g.indx2pt[i]
				p1 := g.indx2pt[j]
				canv.Line(int(p0.X), int(maxY-p0.Y), int(p1.X), int(maxY-p1.Y), "stroke:#DDD000;stroke-width:2")
			}
		}
	}
}

func (g *Graph) DrawInit(canv *svg.SVG) {
	maxX, maxY := g.getMax()
	canv.Start(int(maxX), int(maxY))
	canv.Rect(0, 0, int(maxX), int(maxY), "fill:#CCFFCC")
	for _, plygn := range g.Trrn.Plygns {
		x := make([]int, plygn.N)
		y := make([]int, plygn.N)
		for i, c := range plygn.Corners {
			x[i] = int(c.Pos.X)
			y[i] = int(maxY - c.Pos.Y)
		}
		canv.Polygon(x, y, "fill:#6b6b6b;stroke:black")
	}
}

func (g *Graph) DrawHome(canv *svg.SVG) {
	_, maxY := g.getMax()
	home := g.indx2pt[0]
	x0 := int(home.X)
	y0 := int(maxY - home.Y)
	canv.Circle(x0, y0, 10, "stroke:black;fill:red")
}

func (g *Graph) DrawPath(canv *svg.SVG) {
	x := make([]int, len(g.path)+1)
	y := make([]int, len(g.path)+1)
	_, maxY := g.getMax()
	var p *internal.Point
	for i := 0; i < len(g.path); i++ {
		p = g.indx2pt[g.path[i]]
		x[i] = int(p.X)
		y[i] = int(maxY - p.Y)
	}
	y[len(g.path)] = int(maxY - p.Y - p.X*ry)
	for i := 1; i < len(g.path); i++ {
		canv.Polyline(x, y, "stroke:#000080;stroke-width:3;fill:none")
	}
}
