package internal

import (
	"testing"
)

func TestMakeClockwise(t *testing.T) {
	p := NewPolygon()
	p.AddCorners([]Point{Point{0, 0}, Point{1, 0}, Point{0, 1}})
	p.MakeClockwise()
	if !p.isClockwise() {
		t.FailNow()
	}
	p.Reverse()
	p.MakeClockwise()
	if !p.isClockwise() {
		t.FailNow()
	}
}

func TestReverse(t *testing.T) {
	p := NewPolygon()
	p.AddCorners([]Point{Point{0, 0}, Point{1, 0}, Point{0, 1}})
	b1 := p.isClockwise()
	p.Reverse()
	b2 := p.isClockwise()
	if b1 == b2 {
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
