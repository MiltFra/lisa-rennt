package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/miltfra/lisa-rennt/internal/graph"
	"github.com/miltfra/lisa-rennt/internal/terrain"
)

var g *graph.Graph

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		path = "data/lisarennt5.txt"
	}
	t := terrain.New(path)
	g = graph.New(t)
	g.Shortest()
	http.Handle("/data/", http.HandlerFunc(draw))
	http.Handle("/", http.HandlerFunc(handleFile))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func draw(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	canv := svg.New(w)
	switch strings.Split(req.URL.Path, "/")[2] {
	case "graph":
		g.DrawGraph(canv)
	case "path":
		g.DrawPath(canv)
	case "terrain":
		g.DrawObstacles(canv)
	default:
		g.DrawAll(canv)
	}
	canv.End()
}

func handleFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "data/public/"+r.URL.Path[1:])
}
