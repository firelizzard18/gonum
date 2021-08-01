package rtree

import (
	"fmt"
	"sort"

	"gonum.org/v1/gonum/spatial/curve"
)

func calculateBounds(v BoundedSet) Bounding {
	if v.Len() == 0 {
		return Bounding{}
	}

	b := v.Get(0).Bounds()
	for i, n := 1, v.Len(); i < n; i++ {
		b = b.Extent(v.Get(i).Bounds())
	}
	return b
}

func quantize(b Bounding, w int, p Point) []int {
	v := make([]int, p.Dims())

	for i := range v {
		v[i] = int(float64(w) * (p.Dim(i) - b.Min.Dim(i)) / (b.Max.Dim(i) - b.Min.Dim(i)))
		if v[i] == w {
			v[i]--
		}
	}

	return v
}

type boundedSetSorter struct {
	BoundedSet
	less func(i, j Bounded) bool
}

func (s *boundedSetSorter) Less(i, j int) bool {
	return s.less(s.Get(i), s.Get(j))
}

type hilbertSorter struct {
	set    BoundedSet
	points []int
}

func (s *hilbertSorter) Len() int { return s.set.Len() }

func (s *hilbertSorter) Swap(i, j int) {
	s.set.Swap(i, j)
	s.points[i], s.points[j] = s.points[j], s.points[i]
}

func (s *hilbertSorter) Less(i, j int) bool {
	return s.points[i] < s.points[j]
}

func HilbertCurvePivot(order int) func(BoundedSet) int {
	return func(set BoundedSet) int {
		const k = 5

		if set.Len() < 2 {
			return 0
		}

		var h curve.SpaceFilling
		switch n := set.Get(0).Bounds().Min.Dims(); n {
		case 2:
			h = curve.Hilbert2D{Order: k}
		case 3:
			h = curve.Hilbert3D{Order: k}
		case 4:
			h = curve.Hilbert4D{Order: k}
		default:
			panic(fmt.Errorf("no %d-dimension hilbert curve is not implemented", n))
		}

		bounds := calculateBounds(set)

		sorter := new(hilbertSorter)
		sorter.set = set
		sorter.points = make([]int, set.Len())
		for i, n := 0, set.Len(); i < n; i++ {
			b := set.Get(i)
			bn := b.Bounds()
			center := make(Vector, bn.Min.Dims())
			for i := range center {
				center[i] = (bn.Min.Dim(i) + bn.Max.Dim(i)) / 2
			}

			sorter.points[i] = h.Curve(quantize(bounds, 1<<k, center))
		}

		sort.Sort(sorter)

		return set.Len() / 2
	}
}
