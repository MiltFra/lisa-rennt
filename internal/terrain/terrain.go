package terrain

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/miltfra/lisa-rennt/internal"
	"github.com/miltfra/tools"
)

// Terrain describes a certain instance of the
// Lisa Rennt problem
type Terrain struct {
	PolCnt int
	Plygns []*internal.Polygon
	Start  *internal.Point
}

// New returns a Terrain instance read from a file
func New(path string) *Terrain {
	// Open a stream to the file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("[ERR] Error while reading file:\n", err)
	}
	// Close the stream at the end of this function
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// Read to the next \n
	scanner.Scan()
	polCnt := tools.Stoi(scanner.Text())
	plygns := make([]*internal.Polygon, polCnt)
	for i := 0; i < polCnt; i++ {
		scanner.Scan()
		plygns[i] = internal.SToPlygn(scanner.Text())
	}
	scanner.Scan()
	h := strings.Split(scanner.Text(), " ")
	home := internal.NewPoint(tools.Stof64(h[0]), tools.Stof64(h[1]))
	return &Terrain{polCnt, plygns, home}
}

// GetAccessiblePoints returns a slice of Points which
// should be considered in pathing
func (t *Terrain) GetAccessiblePoints(p *internal.Point) []*internal.Point {
	res := make([]*internal.Point, 0, 10)
	for _, plygn := range t.Plygns {
		res = append(res, plygn.GetUnblockedSpikes(p)...)
	}
	return res
}
