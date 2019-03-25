package graph

import (
	"math"

	"github.com/miltfra/lisa-rennt/internal"
)

// partPath is an instance of the path finding
// algorithm. It represents a sequence of nodes
// starting at 0. It also contains times needed
// to follow the path with velL.
type partPath struct {
	bck    []int
	passed float64
	max    float64
	last   int
}

const segLen = 10000

type prioQ struct {
	data  []*partPath
	len   int
	final []bool
}

// Shortest returns a slice of integers representing
// the indices of nodes in the graph visited before
// walking towards the road in the optimal angle. This
// path is the optimal solution when competing against
// the bus.
func (g *Graph) Shortest() []int {
	bck := make([]int, g.N)
	for i := range bck {
		bck[i] = -1
	}
	start := &partPath{bck, 0, g.maxDelta(g.indx2pt[0], 0), 0}
	q := &prioQ{make([]*partPath, segLen), 0, g.getFinals()}
	q.Put(start)
	var i int
	var part *partPath
	for part = q.Get(); !q.final[part.last]; part = q.Get() {
		for i = 0; i < g.N; i++ {
			if part.bck[i] == -1 && g.Mtrx[part.last*g.N+i] >= 0 {
				q.Put(g.extend(part, i))
			}
		}
	}
	g.path = g.toPath(part)
	g.d = newDrawer(g)
	return g.path
}

func (g *Graph) toPath(p *partPath) []int {
	rev := make([]int, 0, g.N)
	current := p.last
	for {
		rev = append(rev, current)
		if p.bck[current] == -1 || len(rev) == g.N {
			break
		} else {
			current = p.bck[current]
		}
	}
	res := make([]int, len(rev))
	for i := 0; i < len(rev); i++ {
		res[i] = rev[len(rev)-i-1]
	}
	return res
}

func (g *Graph) extend(p *partPath, i int) *partPath {
	newP := &partPath{}
	newP.bck = make([]int, g.N)
	copy(newP.bck, p.bck)
	newP.bck[i] = p.last
	newP.passed = p.passed + g.Mtrx[p.last*g.N+i]/velL
	newP.max = g.maxDelta(g.indx2pt[i], newP.passed)
	newP.last = i
	return newP
}

func (g *Graph) getFinals() []bool {
	res := make([]bool, g.N)
	var optimal internal.Point
	var current *internal.Point
	var ls *internal.LineSegment
	for i := 0; i < g.N; i++ {
		current = g.indx2pt[i]
		optimal = internal.Point{0, ry*current.X + current.Y}
		ls = internal.NewLineSegment(optimal, *current)
		if !HasIntersections(g.Trrn, ls) {
			res[i] = true
		}
	}
	return res
}

const velB = 30 / 3.6
const velL = 15 / 3.6

var ry = math.Sqrt(math.Pow(velB/velL, 2) - 1)

// maxDelta returns the time delta lisa would have when
// travelling from a given point in the optimal route
// without any obstacles. This time is calculated after
// some time has passed already.
func (g *Graph) maxDelta(p *internal.Point, passed float64) float64 {
	// TODO: Optimize calculation (p.X^2 is evalulated more than needed)
	dy := p.X * ry
	ds := math.Sqrt(dy*dy + p.X*p.X)
	return (p.Y+dy)/velB - ds/velL - passed
}

func (q *prioQ) Put(e *partPath) {
	q.len++
	if q.len%segLen == 0 {
		newData := make([]*partPath, q.len+segLen)
		copy(newData, q.data)
		q.data = newData
	}
	q.data[q.len-1] = e
	q.up(q.len - 1)
}

func (q *prioQ) Get() *partPath {
	if q.len == 0 {
		return nil
	}
	v := q.data[0]
	q.len--
	q.data[0] = q.data[q.len]
	q.down(0)
	return v
}

// up moves the element at the given index up until
// the heap condition is achieved for this element.
func (q *prioQ) up(i int) {
	v := q.data[i]
	p := (i - 1) >> 1
	for ; i > 0; p = (i - 1) >> 1 {
		pv := q.data[p]
		if v.greater(pv) {
			q.data[i] = pv
			i = p
		} else {
			break
		}
	}
	q.data[i] = v
}

// down moves the element at the given index down until
// the heap condition is achieved for this element. The
// element is only ever swapped with the bigger child.
func (q *prioQ) down(i int) {
	var cv *partPath
	v := q.data[i]
	c := (i << 1) + 1
	for l := q.len; c+1 < l; c = (i << 1) + 1 {
		if q.data[c].less(q.data[c+1]) {
			c++
		}
		cv = q.data[c]
		if v.less(cv) {
			q.data[i] = cv
			i = c
		} else {
			break
		}
	}
	if c < q.len {
		cv := q.data[c]
		if v.less(cv) {
			q.data[i] = cv
			i = c
		}
	}
	q.data[i] = v
}

func (p *partPath) greater(other *partPath) bool {
	return p.max > other.max
}

func (p *partPath) less(other *partPath) bool {
	return p.max < other.max
}
