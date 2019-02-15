package terrain

import (
	"fmt"
	"testing"
)

func TestMakeClockwise(t *testing.T) {
	trrn := New("/home/miltfra/lisarennttest.txt")
	for i := 0; i < trrn.PolCnt; i++ {
		b1 := trrn.Plygns[i].isClockwise()
		trrn.Plygns[i].Reverse()
		b2 := trrn.Plygns[i].isClockwise()
		trrn.Plygns[i].MakeClockwise()
		b3 := trrn.Plygns[i].isClockwise()
		if !b1 || b2 || !b3 {
			t.FailNow()
		}
	}
}

func TestReverse(t *testing.T) {
	trrn := New("/home/miltfra/lisarennttest.txt")
	for i := 0; i < trrn.PolCnt; i++ {
		b1 := trrn.Plygns[i].isClockwise()
		trrn.Plygns[i].Reverse()
		b2 := trrn.Plygns[i].isClockwise()
		trrn.Plygns[i].Reverse()
		b3 := trrn.Plygns[i].isClockwise()
		if !(b1 != b2 && b1 == b3) {
			t.FailNow()
		}
	}
}

func TestDoIntersect(t *testing.T) {
	ln1 := NewLineSegment(NewPoint(0, 0), NewPoint(0, 1))
	ln2 := NewLineSegment(NewPoint(0, 0), NewPoint(1, 0))
	ln3 := NewLineSegment(NewPoint(-1, 0), NewPoint(1, 1))
	ln4 := NewLineSegment(NewPoint(-1, -0.5), NewPoint(0, 1.5))
	ln5 := NewLineSegment(NewPoint(-1, 0), NewPoint(1, 0))
	ln6 := NewLineSegment(NewPoint(0, 1), NewPoint(0, 0))
	if intersection(ln1, ln2) < 0 {
		t.FailNow()
	}
	if intersection(ln1, ln3) >= 0 {
		t.FailNow()
	}
	if intersection(ln1, ln4) < 0 {
		t.FailNow()
	}
	if intersection(ln1, ln5) < 0 {
		t.FailNow()
	}
	if intersection(ln2, ln3) < 0 {
		t.FailNow()
	}
	if intersection(ln2, ln4) < 0 {
		t.FailNow()
	}
	if intersection(ln2, ln5) < 0 {
		t.FailNow()
	}
	if intersection(ln3, ln4) >= 0 {
		t.FailNow()
	}
	if intersection(ln3, ln5) < 0 {
		t.FailNow()
	}
	if intersection(ln4, ln5) >= 0 {
		t.FailNow()
	}
	if intersection(ln1, ln6) < 0 {
		t.FailNow()
	}
	if intersection(ln6, ln1) < 0 {
		t.FailNow()
	}
	if intersection(ln1, ln1) < 0 {
		t.FailNow()
	}
}

func TestHullContains(t *testing.T) {
	P := NewPolygon()
	P.AddCorners([]Point{{-1, 0}, {0, 0.5}, {1, 0}, {0, 1}})
	if P.HullContains(Point{1, 1}) {
		t.FailNow()
	}
	if !P.HullContains(Point{0, 0.25}) {
		t.FailNow()
	}
	if P.HullContains(Point{0, -.25}) {
		t.FailNow()
	}
	if !P.HullContains(Point{0.25, 0.26}) {
		t.FailNow()
	}
}

func TestDoesIntersect(t *testing.T) {
	p := NewPolygon()
	p.AddCorners([]Point{{-1, 0}, {1, 0}, {0, 1}})
	ls := NewLineSegment(Point{-1, -1}, Point{1, 1})
	if !p.DoesIntersect(ls) {
		t.FailNow()
	}
	ls = NewLineSegment(Point{-1, -1}, Point{1, -1})
	if p.DoesIntersect(ls) {
		t.FailNow()
	}
	ls = NewLineSegment(Point{0, 0}, Point{0, 2})
	if !p.DoesIntersect(ls) {
		t.FailNow()
	}
	ls = NewLineSegment(Point{1, 0}, Point{1, 2})
	if p.DoesIntersect(ls) {
		t.FailNow()
	}
}

func TestGetAngle(t *testing.T) {
	p0 := NewPoint(0, 0)
	a1, err := getAngle(p0, NewPoint(1, 0))
	if err != nil {
		t.FailNow()
	}
	a2, err := getAngle(p0, NewPoint(0, 1))
	if err != nil {
		t.FailNow()
	}
	a3, err := getAngle(p0, NewPoint(-1, 0))
	if err != nil {
		t.FailNow()
	}
	a4, err := getAngle(p0, NewPoint(0, -1))
	if err != nil {
		t.FailNow()
	}
	if a1 >= a2 {
		t.FailNow()
	}
	if a2 >= a3 {
		t.FailNow()
	}
	if a3 >= a4 {
		t.FailNow()
	}
}

//func TestGetShadow(t *testing.T) {
//	p := NewPolygon().AddCorners([]Point{{-1, 0}, {1, 0}, {0, -1}})
//	p1, p2, _ := p.getShadow(Point{0, 1})
//	if p1.Pos != NewPoint(1, 0) && p2.Pos != NewPoint(-1, 0) ||
//		p1.Pos != NewPoint(-1, 0) && p2.Pos != NewPoint(1, 0) {
//		t.FailNow()
//	}
//}

func TestGetAccessiblePoints(t *testing.T) {
	trrn := New("/home/miltfra/lisarennttest.txt")
	pts := GetAccessiblePoints(trrn, Point{0, 0})
	fmt.Println(pts)
}

func TestNewGraph(t *testing.T) {
	trrn := New("/home/miltfra/lisarennt6.txt")
	g := NewGraph(trrn)
	fmt.Println(g.Mtrx)
}

func TestDraw(t *testing.T) {
	trrn := New("/home/miltfra/lisarennt4.txt")
	g := NewGraph(trrn)
	g.Draw("/home/miltfra/lisarennt4.jpg")
}

func TestChangesSide(t *testing.T) {
	ls := NewLineSegment(NewPoint(-1, -1), NewPoint(1, 1))
	if ls.changesSide(NewPoint(0, 0), NewPoint(-1, 0), NewPoint(0, 1)) {
		t.FailNow()
	}
	if ls.changesSide(NewPoint(0, 0), NewPoint(0, 1), NewPoint(-1, 0)) {
		t.FailNow()
	}
	if !ls.changesSide(NewPoint(0, 0), NewPoint(1, 0), NewPoint(0, 1)) {
		t.FailNow()
	}
	if !ls.changesSide(NewPoint(0, 0), NewPoint(0, 1), NewPoint(1, 0)) {
		t.FailNow()
	}
}

func TestIsInCounterClockwiseAngle(t *testing.T) {
	p0 := NewPoint(0, 0)
	p1 := NewPoint(1, 0)
	//p2 := NewPoint(1, 1)
	p3 := NewPoint(0, 1)
	//p4 := NewPoint(-1, 1)
	p5 := NewPoint(-1, 0)
	p6 := NewPoint(0, -1)
	ls := NewLineSegment(p3, p6)
	if ls.inCounterClockwiseAngle(p0, p5, p1) {
		t.FailNow()
	}
	if !ls.inCounterClockwiseAngle(p0, p1, p5) {
		t.FailNow()
	}
	if ls.inCounterClockwiseAngle(p0, p1, p3) {
		t.FailNow()
	}
}
