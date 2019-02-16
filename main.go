package main

import (
	"os"
)

func main() {
	path := os.Args[1]
	t := NewTerrain(path)
	g := NewGraph(t)
	g.Draw(path + ".jpg")
}
