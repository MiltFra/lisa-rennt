package graph

import (
	"fmt"
	"testing"

	"github.com/miltfra/lisa-rennt/internal"

	"github.com/miltfra/lisa-rennt/internal/terrain"
)

func TestGetAccessiblePoints(t *testing.T) {
	trrn := terrain.New("/home/miltfra/lisarennt5.txt")
	pts := trrn.GetAccessiblePoints(internal.NewPoint(0, 0))
	fmt.Println(pts)
}

func TestNewGraph(t *testing.T) {
	trrn := terrain.New("/home/miltfra/lisarennt6.txt")
	g := New(trrn)
	fmt.Println(g.Mtrx)
}

func TestDraw(t *testing.T) {
	trrn := terrain.New("/home/miltfra/lisarennt4.txt")
	g := New(trrn)
	g.Draw("/home/miltfra/lisarennt4.jpg")
}
