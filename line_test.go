package main

import "testing"

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
