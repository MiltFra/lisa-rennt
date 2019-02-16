package main

import "testing"

func TestMakeClockwise(t *testing.T) {
	trrn := NewTerrain("/home/miltfra/lisarennt5.txt")
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
	trrn := NewTerrain("/home/miltfra/lisarennt5.txt")
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
