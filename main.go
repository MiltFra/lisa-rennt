package main

import (
	"log"
	"net/http"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/miltfra/lisa-rennt/internal/graph"
	"github.com/miltfra/lisa-rennt/internal/terrain"
)

func main() {
	http.Handle("/", http.HandlerFunc(draw))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func draw(w http.ResponseWriter, req *http.Request) {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = "data/lisarennt5.txt"
	}
	t := terrain.New(path)
	g := graph.New(t)
	g.Shortest()
	w.Header().Set("Content-Type", "image/svg+xml")
	canv := svg.New(w)
	g.DrawAll(canv)
	canv.End()
}
