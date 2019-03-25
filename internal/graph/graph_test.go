package graph

import (
	"fmt"
	"testing"
)

func TestGetAccessiblePoints(t *testing.T) {
	trrn := NewTerrain("/home/miltfra/lisarennt5.txt")
	pts := GetAccessiblePoints(trrn, Point{0, 0})
	fmt.Println(pts)
}

func TestNewGraph(t *testing.T) {
	trrn := NewTerrain("/home/miltfra/lisarennt6.txt")
	g := NewGraph(trrn)
	fmt.Println(g.Mtrx)
}

func TestDraw(t *testing.T) {
	trrn := NewTerrain("/home/miltfra/lisarennt4.txt")
	g := NewGraph(trrn)
	g.Draw("/home/miltfra/lisarennt4.jpg")
}
