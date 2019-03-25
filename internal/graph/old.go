package graph

//var imgX float64 = 1920
//var imgY float64 = 1000
//var maxX, maxY float64
//
//// Draw draws the graph to an image
//func (G *Graph) Draw(path string) {
//	maxX, maxY = G.getMax()
//	img := image.NewRGBA64(image.Rect(0, 0, int(imgX), int(imgY)))
//	// Colors
//	//fillCol := color.RGBA{100, 100, 100, 255}
//	// Draw Polygons
//	plygnCol := color.RGBA{255, 255, 255, 255}
//	for _, plygn := range G.Trrn.Plygns {
//		//drawPolygon(plygn, img, fillCol)
//		for _, c := range plygn.Corners {
//			x0, y0 := scale(c.Pos.X, c.Pos.Y)
//			x1, y1 := scale(c.N2.Pos.X, c.N2.Pos.Y)
//			drawLine(x0, y0, x1, y1, img, plygnCol)
//		}
//	}
//	if G.path == nil {
//		pathCol := color.RGBA{240, 50, 50, 255}
//		for i := 0; i < G.N; i++ {
//			for j := i + 1; j < G.N; j++ {
//				if G.Mtrx[i*G.N+j] != -1 {
//					x0, y0 := scale(G.indx2pt[i].X, G.indx2pt[i].Y)
//					x1, y1 := scale(G.indx2pt[j].X, G.indx2pt[j].Y)
//					drawLine(x0, y0, x1, y1, img, pathCol)
//				}
//			}
//		}
//	} else {
//		shortestCol := color.RGBA{50, 50, 240, 255}
//		for i := 1; i < len(G.path); i++ {
//			p1 := G.indx2pt[G.path[i-1]]
//			p2 := G.indx2pt[G.path[i]]
//			x0, y0 := scale(p1.X, p1.Y)
//			x1, y1 := scale(p2.X, p2.Y)
//			drawLine(x0, y0, x1, y1, img, shortestCol)
//		}
//		p0 := G.indx2pt[G.path[len(G.path)-1]]
//		p1 := internal.NewPoint(0, p0.Y+p0.X*ry)
//		x0, y0 := scale(p0.X, p0.Y)
//		x1, y1 := scale(p1.X, p1.Y)
//		drawLine(x0, y0, x1, y1, img, shortestCol)
//
//	}
//	toimg, _ := os.Create(path)
//	defer toimg.Close()
//	jpeg.Encode(toimg, img, nil)
//}
//
//func scale(x, y float64) (int, int) {
//	factor := math.Min(imgX/maxX, imgY/maxY)
//	return int(x * factor), int(y * factor)
//}
//
//func abs(x int) int {
//	if x >= 0 {
//		return x
//	}
//	return -x
//}
//
//func drawPolygon(p *internal.Polygon, img *image.RGBA64, col color.Color) {
//	min, max := p.GetBox()
//	for i := min.Y - 1; i <= max.Y+1; i++ {
//		nodes := make([]float64, 0, 20)
//		//ls := NewLineSegment(NewPoint(min.X-1, i), NewPoint(max.X+1, i))
//		for _, c := range p.Corners {
//			if (c.Pos.Y <= i) != (c.N2.Pos.Y <= i) {
//				nodes = append(nodes, c.Pos.X+(i-c.Pos.Y)/(c.N2.Pos.Y-c.Pos.Y)*(c.N2.Pos.X-c.Pos.X))
//			}
//		}
//		sort.Float64s(nodes)
//		for k := 0; k < len(nodes)-1; k += 2 {
//			x0, y0 := scale(nodes[k], i)
//			x1, y1 := scale(nodes[k+1], i)
//			drawLine(x0, y0, x1, y1, img, col)
//		}
//	}
//}
//
//// drawLine draws a line between two points on an image
//func drawLine(x0, y0, x1, y1 int, img *image.RGBA64, col color.Color) {
//	dx := abs(x1 - x0)
//	dy := abs(y1 - y0)
//	sx, sy := 1, 1
//	if x0 >= x1 {
//		sx = -1
//	}
//	if y0 >= y1 {
//		sy = -1
//	}
//	err := dx - dy
//	for {
//		img.Set(x0, y0, col)
//		if x0 == x1 && y0 == y1 {
//			return
//		}
//		e2 := err * 2
//		if e2 > -dy {
//			err -= dy
//			x0 += sx
//		}
//		if e2 < dx {
//			err += dx
//			y0 += sy
//		}
//	}
//}
