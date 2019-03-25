package internal

import "testing"

func TestGetAngle(t *testing.T) {
	p0 := NewPoint(0, 0)
	a1, err := GetAngle(p0, NewPoint(1, 0))
	if err != nil {
		t.FailNow()
	}
	a2, err := GetAngle(p0, NewPoint(0, 1))
	if err != nil {
		t.FailNow()
	}
	a3, err := GetAngle(p0, NewPoint(-1, 0))
	if err != nil {
		t.FailNow()
	}
	a4, err := GetAngle(p0, NewPoint(0, -1))
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
