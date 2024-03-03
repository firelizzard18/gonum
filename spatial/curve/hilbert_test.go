// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curve

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/exp/rand"
)

func ExampleHilbert2D_Pos() {
	h := Hilbert2D{Order: 3}

	for y := 0; y < 1<<h.Order; y++ {
		for x := 0; x < 1<<h.Order; x++ {
			if x > 0 {
				fmt.Print("  ")
			}
			fmt.Printf("%02X", h.Pos([]int{x, y}))
		}
		fmt.Println()
	}
	// Output:
	// 00  01  0E  0F  10  13  14  15
	// 03  02  0D  0C  11  12  17  16
	// 04  07  08  0B  1E  1D  18  19
	// 05  06  09  0A  1F  1C  1B  1A
	// 3A  39  36  35  20  23  24  25
	// 3B  38  37  34  21  22  27  26
	// 3C  3D  32  33  2E  2D  28  29
	// 3F  3E  31  30  2F  2C  2B  2A
}

func ExampleHilbert3D_Pos() {
	h := Hilbert3D{Order: 2}

	for z := 0; z < 1<<h.Order; z++ {
		for y := 0; y < 1<<h.Order; y++ {
			for x := 0; x < 1<<h.Order; x++ {
				if x > 0 {
					fmt.Print("  ")
				}
				fmt.Printf("%02X", h.Pos([]int{x, y, z}))
			}
			fmt.Println()
		}
		fmt.Println()
	}
	// Output:
	// 00  07  08  0B
	// 01  06  0F  0C
	// 1A  1B  10  13
	// 19  18  17  14
	//
	// 03  04  09  0A
	// 02  05  0E  0D
	// 1D  1C  11  12
	// 1E  1F  16  15
	//
	// 3C  3B  36  35
	// 3D  3A  31  32
	// 22  23  2E  2D
	// 21  20  29  2A
	//
	// 3F  38  37  34
	// 3E  39  30  33
	// 25  24  2F  2C
	// 26  27  28  2B
}

func ExampleHilbert4D_Pos() {
	h := Hilbert4D{Order: 2}
	N := 1 << h.Order
	for z := 0; z < N; z++ {
		if z > 0 {
			s := strings.Repeat("═", N*4-2)
			s = s + strings.Repeat("═╬═"+s, N-1)
			fmt.Println(s)
		}
		for y := 0; y < N; y++ {
			for w := 0; w < N; w++ {
				if w > 0 {
					fmt.Print(" ║ ")
				}
				for x := 0; x < N; x++ {
					if x > 0 {
						fmt.Print("  ")
					}
					fmt.Printf("%02X", h.Pos([]int{x, y, z, w}))
				}
			}
			fmt.Println()
		}
	}
	// Output:
	// 00  0F  10  13 ║ 03  0C  11  12 ║ FC  F3  EE  ED ║ FF  F0  EF  EC
	// 01  0E  1F  1C ║ 02  0D  1E  1D ║ FD  F2  E1  E2 ║ FE  F1  E0  E3
	// 32  31  20  23 ║ 35  36  21  22 ║ CA  C9  DE  DD ║ CD  CE  DF  DC
	// 33  30  2F  2C ║ 34  37  2E  2D ║ CB  C8  D1  D2 ║ CC  CF  D0  D3
	// ═══════════════╬════════════════╬════════════════╬═══════════════
	// 07  08  17  14 ║ 04  0B  16  15 ║ FB  F4  E9  EA ║ F8  F7  E8  EB
	// 06  09  18  1B ║ 05  0A  19  1A ║ FA  F5  E6  E5 ║ F9  F6  E7  E4
	// 3D  3E  27  24 ║ 3A  39  26  25 ║ C5  C6  D9  DA ║ C2  C1  D8  DB
	// 3C  3F  28  2B ║ 3B  38  29  2A ║ C4  C7  D6  D5 ║ C3  C0  D7  D4
	// ═══════════════╬════════════════╬════════════════╬═══════════════
	// 76  77  6C  6D ║ 79  78  6B  6A ║ 86  87  94  95 ║ 89  88  93  92
	// 75  74  63  62 ║ 7A  7B  64  65 ║ 85  84  9B  9A ║ 8A  8B  9C  9D
	// 42  41  5C  5D ║ 45  46  5B  5A ║ BA  B9  A4  A5 ║ BD  BE  A3  A2
	// 43  40  53  52 ║ 44  47  54  55 ║ BB  B8  AB  AA ║ BC  BF  AC  AD
	// ═══════════════╬════════════════╬════════════════╬═══════════════
	// 71  70  6F  6E ║ 7E  7F  68  69 ║ 81  80  97  96 ║ 8E  8F  90  91
	// 72  73  60  61 ║ 7D  7C  67  66 ║ 82  83  98  99 ║ 8D  8C  9F  9E
	// 4D  4E  5F  5E ║ 4A  49  58  59 ║ B5  B6  A7  A6 ║ B2  B1  A0  A1
	// 4C  4F  50  51 ║ 4B  48  57  56 ║ B4  B7  A8  A9 ║ B3  B0  AF  AE
}

func TestHilbert2D(t *testing.T) {
	for ord := 1; ord <= 4; ord++ {
		t.Run(fmt.Sprintf("Order/%d", ord), func(t *testing.T) {
			testCurve(t, Hilbert2D{Order: ord})
		})
	}

	cases := map[int][]int{
		1: {
			0, 1,
			3, 2,
		},
		2: {
			0x0, 0x3, 0x4, 0x5,
			0x1, 0x2, 0x7, 0x6,
			0xE, 0xD, 0x8, 0x9,
			0xF, 0xC, 0xB, 0xA,
		},
		3: {
			0x00, 0x01, 0x0E, 0x0F, 0x10, 0x13, 0x14, 0x15,
			0x03, 0x02, 0x0D, 0x0C, 0x11, 0x12, 0x17, 0x16,
			0x04, 0x07, 0x08, 0x0B, 0x1E, 0x1D, 0x18, 0x19,
			0x05, 0x06, 0x09, 0x0A, 0x1F, 0x1C, 0x1B, 0x1A,
			0x3A, 0x39, 0x36, 0x35, 0x20, 0x23, 0x24, 0x25,
			0x3B, 0x38, 0x37, 0x34, 0x21, 0x22, 0x27, 0x26,
			0x3C, 0x3D, 0x32, 0x33, 0x2E, 0x2D, 0x28, 0x29,
			0x3F, 0x3E, 0x31, 0x30, 0x2F, 0x2C, 0x2B, 0x2A,
		},
	}

	for order, expected := range cases {
		t.Run(fmt.Sprintf("Case/%d", order), func(t *testing.T) {
			testCurveCase(t, Hilbert2D{Order: order}, order, expected)
		})
	}
}

func TestHilbert3D(t *testing.T) {
	for ord := 1; ord <= 4; ord++ {
		t.Run(fmt.Sprintf("Order/%d", ord), func(t *testing.T) {
			testCurve(t, Hilbert3D{Order: ord})
		})
	}

	cases := map[int][]int{
		1: {
			0, 1,
			3, 2,

			7, 6,
			4, 5,
		},
		2: {
			0x00, 0x07, 0x08, 0x0B,
			0x01, 0x06, 0x0F, 0x0C,
			0x1A, 0x1B, 0x10, 0x13,
			0x19, 0x18, 0x17, 0x14,

			0x03, 0x04, 0x09, 0x0A,
			0x02, 0x05, 0x0E, 0x0D,
			0x1D, 0x1C, 0x11, 0x12,
			0x1E, 0x1F, 0x16, 0x15,

			0x3C, 0x3B, 0x36, 0x35,
			0x3D, 0x3A, 0x31, 0x32,
			0x22, 0x23, 0x2E, 0x2D,
			0x21, 0x20, 0x29, 0x2A,

			0x3F, 0x38, 0x37, 0x34,
			0x3E, 0x39, 0x30, 0x33,
			0x25, 0x24, 0x2F, 0x2C,
			0x26, 0x27, 0x28, 0x2B,
		},
	}

	for order, expected := range cases {
		t.Run(fmt.Sprintf("Case/%d", order), func(t *testing.T) {
			testCurveCase(t, Hilbert3D{Order: order}, order, expected)
		})
	}
}

func TestHilbert4D(t *testing.T) {
	for ord := 1; ord <= 4; ord++ {
		t.Run(fmt.Sprintf("Order %d", ord), func(t *testing.T) {
			testCurve(t, Hilbert4D{Order: ord})
		})
	}
}

func BenchmarkHilbert(b *testing.B) {
	const O = 10
	for N := 2; N <= 4; N++ {
		b.Run(fmt.Sprintf("%dD/Pos", N), func(b *testing.B) {
			for ord := 1; ord <= O; ord++ {
				b.Run(fmt.Sprintf("Order %d", ord), func(b *testing.B) {
					h := newCurve(ord, N)
					v := make([]int, N)
					for i := range v {
						v[i] = rand.Intn(1 << ord)
					}
					u := make([]int, N)
					for n := 0; n < b.N; n++ {
						copy(u, v)
						h.Pos(u)
					}
				})
			}
		})

		b.Run(fmt.Sprintf("%dD/Coord", N), func(b *testing.B) {
			for ord := 1; ord <= O; ord++ {
				b.Run(fmt.Sprintf("Order %d", ord), func(b *testing.B) {
					h := newCurve(ord, N)
					d := rand.Intn(1 << (ord * N))
					for n := 0; n < b.N; n++ {
						h.Coord(d)
					}
				})
			}
		})
	}
}

type curve interface {
	Dims() []int
	Len() int
	Pos(v []int) int
	Coord(pos int) []int
}

func newCurve(order, dim int) curve {
	switch dim {
	case 2:
		return Hilbert2D{Order: order}
	case 3:
		return Hilbert3D{Order: order}
	case 4:
		return Hilbert4D{Order: order}
	}
	panic("invalid dimension")
}

func testCurve(t *testing.T, c curve) {
	t.Helper()

	var errc int
	fail := func() {
		if errc < 10 {
			errc++
			t.Fail()
		} else {
			t.FailNow()
		}
	}

	m := map[int][]int{}
	curveRange(c, func(v []int) {
		d := c.Pos(dup(v))
		u := c.Coord(d)
		if !reflect.DeepEqual(v, u) {
			t.Logf("Space is not the inverse of Curve for d=%d %v", d, v)
			fail()
		}

		m[d] = dup(v)
	})

	D := 1
	for _, v := range c.Dims() {
		D *= v
	}
	for d := 0; d < D-1; d++ {
		v, u := m[d], m[d+1]
		if !adjacent(v, u) {
			t.Logf("points %x and %x are not adjacent", d, d+1)
			t.Logf("    %v -> %v", v, u)
			fail()
		}
	}
}

func curveRange(c curve, fn func([]int)) {
	size := c.Dims()
	dimRange(len(size), size, make([]int, len(size)), fn)
}

func dimRange(dim int, size []int, v []int, fn func([]int)) {
	if dim == 0 {
		fn(v)
		return
	}

	for i := 0; i < size[dim-1]; i++ {
		v[dim-1] = i
		dimRange(dim-1, size, v, fn)
	}
}

// v may get mangled by Curve, so copy it.
func dup(v []int) []int {
	u := make([]int, len(v))
	copy(u, v)
	return u
}

// adjacent returns true if the euclidean distance between v and u is
// exactly one. v and u must be the same length.
//
// In other words, given d(i) = abs(v[i] - u[i]), adjacent returns false if
// d(i) > 1 for any i or if d(i) == 1 for more than a single i, or if d(i)
// == 0 for all i.
func adjacent(v, u []int) bool {
	n := 0
	for i := range v {
		x := v[i] - u[i]
		if x == 0 {
			continue
		}
		if x < -1 || 1 < x {
			return false
		}
		n++
	}

	return n == 1
}

func testCurveCase(t *testing.T, c curve, order int, want []int) {
	t.Helper()

	dim := len(c.Dims())
	got := make([]int, len(want))
	for i := range want {
		v := coord(i, order, dim)
		got[i] = c.Pos(dup(v))
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("unexpected result:\ngot:  %v\nwant: %v", got, want)
	}

	for i, d := range want {
		v := coord(i, order, dim)
		if !reflect.DeepEqual(v, c.Coord(d)) {
			t.Fatalf("[%d] expected %v for d = %d", i, v, d)
		}
	}
}

// coord returns the nth spatial coordinates for a dim-dimensional space with
// sides 2^ord in row-major order.
func coord(n, ord, dim int) []int {
	v := make([]int, dim)
	for i := 0; i < dim; i++ {
		v[i] = n % (1 << ord)
		n /= (1 << ord)
	}
	return v
}
