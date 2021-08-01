package rtree_test

import (
	"fmt"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/spatial/rtree"
)

type Vector []float64

func (v Vector) Dims() int         { return len(v) }
func (v Vector) Dim(i int) float64 { return v[i] }

func (v Vector) String() string {
	var s = "<"
	for i, v := range v {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%.2f", v)
	}
	return s + ">"
}

type Cloud []Vector

func (c Cloud) Bounds() rtree.Bounding {
	b := rtree.Bounding{c[0], c[0]}
	for _, v := range c[1:] {
		b = b.Extent(rtree.Bounding{v, v})
	}
	return b
}

type CloudBounded struct {
	c Cloud
	b rtree.Bounding
}

func (c *CloudBounded) Bounds() rtree.Bounding { return c.b }

func dumpTree(n rtree.Node, i0, indent string) {
	switch n := n.(type) {
	case *rtree.Leaf:
		if len(n.Values) == 1 {
			fmt.Printf("%s─%v\n", i0, n.Values[0])
			break
		}
		for i, v := range n.Values {
			switch i {
			case 0:
				fmt.Printf("%s┬%v\n", i0, v)
			case len(n.Values) - 1:
				fmt.Printf("%s└%v\n", indent, v)
			default:
				fmt.Printf("%s├%v\n", indent, v)
			}
		}
	case *rtree.Branch:
		if len(n.Children) == 1 {
			dumpTree(n.Children[0], i0+"─", indent+" ")
			break
		}
		for i, c := range n.Children {
			switch i {
			case 0:
				dumpTree(c, i0+"┬", indent+"|")
			case len(n.Children) - 1:
				dumpTree(c, indent+"└", indent+" ")
			default:
				dumpTree(c, indent+"├", indent+"|")
			}
		}
	}
}

func normalPoint(mean, stddev float64) float64 {
	return rand.NormFloat64()*stddev + mean
}

func ExampleTree() {
	const stddev = 0.1

	centers := []rtree.Vector{
		{0.5, 0.5, 0.5},
		{2.5, 0.5, 0.5},
		{0.5, 2.5, 0.5},
		{0.5, 0.5, 2.5},
		{0.5, 2.5, 2.5},
		{2.5, 0.5, 2.5},
		{2.5, 2.5, 0.5},
		{2.5, 2.5, 2.5},
	}

	tree := new(rtree.Tree)
	tree.FillLevel = 2

	for i := 0; i < len(centers)*3; i++ {
		center := centers[rand.Intn(len(centers))]

		cloud := make(Cloud, 1)
		for i := range cloud {
			cloud[i] = Vector{
				normalPoint(center[0], stddev),
				normalPoint(center[1], stddev),
				normalPoint(center[2], stddev),
			}
		}

		tree.Insert(cloud)
	}

	// Output:
	// ┬┬┬[<2.44, 2.37, 0.61>]
	// ||└[<0.56, 2.62, 0.69>]
	// |├┬[<2.40, 2.47, 0.40>]
	// ||├[<2.43, 2.53, 0.46>]
	// ||├[<2.41, 0.40, 0.55>]
	// ||└[<2.46, 0.56, 0.61>]
	// |├┬[<2.50, 2.39, 0.55>]
	// ||├[<2.42, 0.57, 0.54>]
	// ||└[<0.54, 0.57, 0.48>]
	// |└┬[<2.51, 2.49, 2.37>]
	// | ├[<2.37, 2.48, 2.33>]
	// | ├[<0.48, 0.57, 2.34>]
	// | └[<0.45, 0.46, 2.32>]
	// └┬┬[<0.56, 2.41, 2.46>]
	//  |├[<2.54, 2.68, 2.57>]
	//  |├[<2.58, 0.42, 2.53>]
	//  |└[<2.48, 0.53, 2.49>]
	//  ├┬[<2.59, 0.52, 2.46>]
	//  |├[<0.60, 0.66, 2.47>]
	//  |├[<0.51, 0.43, 2.49>]
	//  |└[<2.54, 0.45, 2.59>]
	//  └┬[<0.38, 0.45, 2.42>]
	//   ├[<2.44, 2.55, 2.46>]
	//   └[<0.57, 0.47, 2.44>]
	dumpTree(tree.Root, "", "")
}

func TestTreeInsert(t *testing.T) {
	tree := new(rtree.Tree)
	tree.FillLevel = 1

	tree.Insert(Cloud{{0, 0, 0}, {1, 1, 1}})
	tree.Insert(Cloud{{2, 0, 0}, {3, 1, 1}})
	tree.Insert(Cloud{{0, 2, 0}, {1, 3, 1}})
	tree.Insert(Cloud{{0, 0, 2}, {1, 1, 3}})
	tree.Insert(Cloud{{2, 2, 0}, {3, 3, 1}})
	tree.Insert(Cloud{{0, 2, 2}, {1, 3, 3}})
	tree.Insert(Cloud{{2, 0, 2}, {3, 1, 3}})
	tree.Insert(Cloud{{2, 2, 2}, {3, 3, 3}})

	if tree.Len() != 8 {
		t.Fatalf("Tree length:\ngot:  %d\nwant: %d\n", tree.Len(), 8)
	}

	bounds := Cloud{{0, 0, 0}, {3, 3, 3}}
	if !tree.Root.Bounds().EqualTo(bounds.Bounds()) {
		t.Fatalf("Tree length:\ngot:  %#v\nwant: %#v\n", tree.Root.Bounds(), bounds.Bounds())
	}
}

func TestTreeSearch(t *testing.T) {
	cases := []struct {
		count  int
		bounds Cloud
	}{
		{8, Cloud{{0, 0, 0}, {3, 3, 3}}},
		{8, Cloud{{1, 1, 1}, {2, 2, 2}}},
		{0, Cloud{{1.1, 1.1, 1.1}, {1.9, 1.9, 1.9}}},
		{1, Cloud{{0, 0, 0}, {1, 1, 1}}},
		{2, Cloud{{0, 0, 0}, {3, 1, 1}}},
		{4, Cloud{{0, 0, 0}, {3, 3, 1}}},
	}

	tree := new(rtree.Tree)
	tree.FillLevel = 1

	tree.Insert(Cloud{{0, 0, 0}, {1, 1, 1}})
	tree.Insert(Cloud{{2, 0, 0}, {3, 1, 1}})
	tree.Insert(Cloud{{0, 2, 0}, {1, 3, 1}})
	tree.Insert(Cloud{{0, 0, 2}, {1, 1, 3}})
	tree.Insert(Cloud{{2, 2, 0}, {3, 3, 1}})
	tree.Insert(Cloud{{0, 2, 2}, {1, 3, 3}})
	tree.Insert(Cloud{{2, 0, 2}, {3, 1, 3}})
	tree.Insert(Cloud{{2, 2, 2}, {3, 3, 3}})

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			k := rtree.NewOverlap(c.bounds)
			tree.Search(k)
			if len(k.Items) != c.count {
				t.Fatalf("Results:\ngot:  %d\nwant: %d\n", len(k.Items), c.count)
			}
		})
	}
}

func BenchmarkTreeSearch(b *testing.B) {
	tree := new(rtree.Tree)
	tree.FillLevel = 5
	for i := 0; i < 1000; i++ {
		v := make(Vector, 3)
		for i := range v {
			v[i] = rand.NormFloat64()
		}

		tree.Insert(&CloudBounded{c: Cloud{v}, b: Cloud{v}.Bounds()})
	}

	search := Cloud{{-0.1, -0.1, -0.1}, {+0.1, +0.1, +0.1}}
	k := rtree.NewOverlap(&CloudBounded{c: search, b: search.Bounds()})
	k.Items = make([]rtree.Bounded, 10)

	for i := 0; i < b.N; i++ {
		k.Items = k.Items[:0]
		tree.Search(k)
	}
}
