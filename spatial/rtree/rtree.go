package rtree

import (
	"math"
	"sort"
)

type Point interface {
	Dims() int
	Dim(int) float64
}

type Vector []float64

func (v Vector) Dims() int         { return len(v) }
func (v Vector) Dim(i int) float64 { return v[i] }

type Bounded interface {
	Bounds() Bounding
}

type Bounding struct {
	Min, Max Point
}

// Overlaps returns true if A and B overlap each other.
func (a Bounding) Overlaps(b Bounding) bool {
	// A overlaps B if and only if a_min_d is less than b_max_d and a_max_d is
	// greater than b_min_d, for each dimension d
	for i, n := 0, a.Min.Dims(); i < n; i++ {
		if a.Min.Dim(i) > b.Max.Dim(i) || a.Max.Dim(i) < b.Min.Dim(i) {
			return false
		}
	}

	return true
}

// Extent calculates the minimum bounds containing both A and B.
func (a Bounding) Extent(b Bounding) Bounding {
	n := a.Min.Dims()
	min, max := make(Vector, n), make(Vector, n)
	for i := 0; i < n; i++ {
		min[i] = math.Min(a.Min.Dim(i), b.Min.Dim(i))
		max[i] = math.Max(a.Max.Dim(i), b.Max.Dim(i))
	}
	return Bounding{Min: min, Max: max}
}

// Area calculates the area of the bounding box.
func (a Bounding) Area() float64 {
	var area float64 = 1
	for i, n := 0, a.Min.Dims(); i < n; i++ {
		area *= a.Max.Dim(i) - a.Min.Dim(i)
	}
	return area
}

func (a Bounding) EqualTo(b Bounding) bool {
	for i, n := 0, a.Min.Dims(); i < n; i++ {
		if a.Min.Dim(i) != b.Min.Dim(i) || a.Max.Dim(i) != b.Max.Dim(i) {
			return false
		}
	}
	return true
}

type Keeper interface {
	Include(Bounded) bool
	Keep(Bounded)
}

type OverlapKeeper struct {
	Bounded
	Items []Bounded
}

func NewOverlap(b Bounded) *OverlapKeeper {
	return &OverlapKeeper{Bounded: b}
}

func (k *OverlapKeeper) Include(b Bounded) bool {
	return k.Bounds().Overlaps(b.Bounds())
}

func (k *OverlapKeeper) Keep(b Bounded) {
	if k.Include(b) {
		k.Items = append(k.Items, b)
	}
}

type BoundedSet interface {
	Len() int
	Get(int) Bounded
	Swap(i, j int)
}

type boundedArray []Bounded

func (a boundedArray) Len() int          { return len(a) }
func (a boundedArray) Get(i int) Bounded { return a[i] }
func (a boundedArray) Swap(i, j int)     { a[i], a[j] = a[j], a[i] }

type nodeArray []Node

func (a nodeArray) Len() int          { return len(a) }
func (a nodeArray) Get(i int) Bounded { return a[i] }
func (a nodeArray) Swap(i, j int)     { a[i], a[j] = a[j], a[i] }

type Options struct {
	// FillLevel determines the target fill level of nodes. A node that drops
	// below the fill level will be disolved, and a node that rises above double
	// the fill level will be split.
	FillLevel int

	// Pivot calculates the pivot used to partition the set.
	Pivot func(BoundedSet) int
}

type Node interface {
	Bounded
	len() int
	rebound(Bounded)
	search(Keeper)
	insert(*Options, Bounded) Node
}

type Branch struct {
	Bounding
	Children []Node
}

func (r *Branch) Bounds() Bounding { return r.Bounding }

func (r *Branch) len() int {
	var n int
	for _, c := range r.Children {
		n += c.len()
	}
	return n
}

func (r *Branch) rebound(b Bounded) {
	if b == nil {
		r.Bounding = calculateBounds(nodeArray(r.Children))
	} else {
		r.Bounding = r.Bounding.Extent(b.Bounds())
	}
}

func (r *Branch) search(k Keeper) {
	if !k.Include(r) {
		return
	}

	for _, c := range r.Children {
		c.search(k)
	}
}

func (r *Branch) pickForInsert(b Bounded) Node {
	candidates := make([]struct {
		node Node
		cost float64
	}, len(r.Children))

	bbounds := b.Bounds()
	for i, c := range r.Children {
		cbounds := c.Bounds()
		cost := cbounds.Extent(bbounds).Area() - cbounds.Area()
		if cost == 0 {
			return c
		}

		candidates[i].node = c
		candidates[i].cost = cost
	}

	sort.Slice(candidates, func(i, j int) bool { return candidates[i].cost < candidates[j].cost })

	return candidates[0].node
}

func (r *Branch) insert(o *Options, b Bounded) Node {
	picked := r.pickForInsert(b)
	split := picked.insert(o, b)
	if split == nil {
		r.rebound(picked)
		return nil
	}

	r.Children = append(r.Children, split)
	if len(r.Children) <= o.FillLevel*2 {
		return nil
	}

	s := new(Branch)
	i := o.Pivot(nodeArray(r.Children))
	s.Children = make([]Node, 0, o.FillLevel*2)
	s.Children = append(s.Children, r.Children[i:]...)
	r.Children = r.Children[:i]

	r.rebound(nil)
	s.rebound(nil)
	return s
}

type Leaf struct {
	Bounding
	Values []Bounded
}

func (l *Leaf) rebound(b Bounded) {
	if b == nil {
		l.Bounding = calculateBounds(boundedArray(l.Values))
	} else {
		l.Bounding = l.Bounding.Extent(b.Bounds())
	}
}

func (l *Leaf) Bounds() Bounding { return l.Bounding }

func (l *Leaf) len() int { return len(l.Values) }

func (l *Leaf) search(k Keeper) {
	if !k.Include(l) {
		return
	}

	for _, c := range l.Values {
		k.Keep(c)
	}
}

func (l *Leaf) insert(o *Options, b Bounded) Node {
	l.Values = append(l.Values, b)
	if len(l.Values) <= o.FillLevel*2 {
		l.rebound(b)
		return nil
	}

	s := new(Leaf)
	i := o.Pivot(boundedArray(l.Values))
	s.Values = make([]Bounded, 0, o.FillLevel*2)
	s.Values = append(s.Values, l.Values[i:]...)
	l.Values = l.Values[:i]

	l.rebound(nil)
	s.rebound(nil)
	return s
}

var DefaultOptions = Options{
	FillLevel: 2,
	Pivot:     HilbertCurvePivot(5),
}

type Tree struct {
	Options
	Root Node
}

func (t *Tree) Len() int { return t.Root.len() }

func (t *Tree) Search(k Keeper) {
	if t.Root != nil {
		t.Root.search(k)
	}
}

func (t *Tree) Insert(b Bounded) {
	if t.FillLevel < 1 {
		t.FillLevel = DefaultOptions.FillLevel
	}
	if t.Pivot == nil {
		n := b.Bounds().Min.Dims()
		if 2 <= n || n <= 4 {
			t.Pivot = DefaultOptions.Pivot
		} else {
			panic("pivot not set")
		}
	}

	if t.Root == nil {
		l := new(Leaf)
		l.Values = make([]Bounded, 0, t.Options.FillLevel*2)
		l.Values = append(l.Values, b)
		l.rebound(nil)
		t.Root = l
		return
	}

	split := t.Root.insert(&t.Options, b)
	if split == nil {
		return
	}

	r := new(Branch)
	r.Children = make([]Node, 0, t.FillLevel*2)
	r.Children = append(r.Children, t.Root, split)
	r.rebound(nil)
	t.Root = r
}
