package main

import (
	"os"

	"github.com/miltfra/lisa-rennt/terrain"
)

func main() {
	path := os.Args[1]
	t := terrain.New(path)
	g := terrain.NewGraph(t)
	g.Draw(path + ".jpg")
}
