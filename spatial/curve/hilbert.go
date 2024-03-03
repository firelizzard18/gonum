// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curve

// Hilbert2D is a 2-dimensional Hilbert curve.
type Hilbert2D struct{ Order int }

// Dims returns the spatial dimensions of the curve, which is {2ᵏ, 2ᵏ}, where k
// is the order.
func (h Hilbert2D) Dims() []int { return []int{1 << h.Order, 1 << h.Order} }

// Len returns the length of the curve, which is 2ⁿᵏ, where n is the dimension
// (2) and k is the order.
//
// Len will overflow on a 32-bit architecture if the order is ≥ 16.
func (h Hilbert2D) Len() int { return 1 << (2 * h.Order) }

func (h Hilbert2D) rot(n int, v []int, d int) {
	switch d {
	case 0:
		swap{0, 1}.do(n, v)
	case 3:
		flip{0, 1}.do(n, v)
	}
}

// Pos returns the linear position of the spatial coordinate along the curve.
// Pos modifies v.
//
// Pos will overflow on a 32-bit architecture if the order is ≥ 16.
func (h Hilbert2D) Pos(v []int) int {
	var d int
	for n := h.Order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rd := ry<<1 | (ry ^ rx)
		d += rd << (2 * n)
		h.rot(h.Order, v, rd)
	}
	return d
}

// Coord2D returns the spatial coordinates of pos.
func (h Hilbert2D) Coord2D(pos int) [2]int {
	var v [2]int
	for n := 0; n < h.Order; n++ {
		e := pos & 3
		h.rot(n, v[:], e)

		ry := e >> 1
		rx := (e>>0 ^ e>>1) & 1
		v[0] += rx << n
		v[1] += ry << n
		pos >>= 2
	}
	return v
}

// Coord returns the spatial coordinates of pos as a slice.
func (h Hilbert2D) Coord(pos int) []int {
	v := h.Coord2D(pos)
	return v[:]
}

// Hilbert3D is a 3-dimensional Hilbert curve.
type Hilbert3D struct{ Order int }

// Dims returns the spatial dimensions of the curve, which is {2ᵏ, 2ᵏ, 2ᵏ}, where
// k is the order.
func (h Hilbert3D) Dims() []int { return []int{1 << h.Order, 1 << h.Order, 1 << h.Order} }

// Len returns the length of the curve, which is 2ⁿᵏ, where n is the dimension
// (3) and k is the order.
//
// Len will overflow on a 32-bit architecture if the order is ≥ 11.
func (h Hilbert3D) Len() int { return 1 << (3 * h.Order) }

func (h Hilbert3D) rot(reverse bool, n int, v []int, d int) {
	switch d {
	case 0:
		do2(reverse, n, v, swap{1, 2}, swap{0, 2})
	case 1, 2:
		do2(reverse, n, v, swap{0, 2}, swap{1, 2})
	case 3, 4:
		invert{0, 1}.do(n, v)
	case 5, 6:
		do2(reverse, n, v, flip{0, 2}, flip{1, 2})
	case 7:
		do2(reverse, n, v, flip{1, 2}, flip{0, 2})
	}
}

// Pos returns the linear position of the spatial coordinate along the curve.
// Pos modifies v.
//
// Pos will overflow on a 32-bit architecture if the order is ≥ 11.
func (h Hilbert3D) Pos(v []int) int {
	var d int
	for n := h.Order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rz := (v[2] >> n) & 1
		rd := rz<<2 | (rz^ry)<<1 | (rz ^ ry ^ rx)
		d += rd << (3 * n)
		h.rot(false, h.Order, v, rd)
	}
	return d
}

// Coord3D returns the spatial coordinates of pos.
func (h Hilbert3D) Coord3D(pos int) [3]int {
	var v [3]int
	for n := 0; n < h.Order; n++ {
		e := pos & 7
		h.rot(true, n, v[:], e)

		rz := e >> 2
		ry := (e>>1 ^ e>>2) & 1
		rx := (e>>0 ^ e>>1) & 1
		v[0] += rx << n
		v[1] += ry << n
		v[2] += rz << n
		pos >>= 3
	}
	return v
}

// Coord returns the spatial coordinates of pos as a slice.
func (h Hilbert3D) Coord(pos int) []int {
	v := h.Coord3D(pos)
	return v[:]
}

// Hilbert4D is a 4-dimensional Hilbert curve.
type Hilbert4D struct{ Order int }

// Dims returns the spatial dimensions of the curve, which is {2ᵏ, 2ᵏ, 2ᵏ, 2ᵏ},
// where k is the order.
func (h Hilbert4D) Dims() []int { return []int{1 << h.Order, 1 << h.Order, 1 << h.Order, 1 << h.Order} }

// Len returns the length of the curve, which is 2ⁿᵏ, where n is the dimension
// (4) and k is the order.
//
// Len will overflow on a 32-bit architecture if the order is ≥ 8.
func (h Hilbert4D) Len() int { return 1 << (4 * h.Order) }

func (h Hilbert4D) rot(reverse bool, n int, v []int, d int) {
	switch d {
	case 0:
		do2(reverse, n, v, swap{1, 3}, swap{0, 3})
	case 1, 2:
		do2(reverse, n, v, swap{0, 3}, swap{1, 3})
	case 3, 4:
		do2(reverse, n, v, flip{0, 1}, swap{2, 3})
	case 5, 6:
		do2(reverse, n, v, flip{1, 2}, swap{2, 3})
	case 7, 8:
		invert{0, 2}.do(n, v)
	case 9, 10:
		do2(reverse, n, v, flip{1, 2}, flip{2, 3})
	case 11, 12:
		do2(reverse, n, v, flip{0, 1}, flip{2, 3})
	case 13, 14:
		do2(reverse, n, v, flip{0, 3}, flip{1, 3})
	case 15:
		do2(reverse, n, v, flip{1, 3}, flip{0, 3})
	}
}

// Pos returns the linear position of the spatial coordinate along the curve.
// Pos modifies v.
//
// Pos will overflow on a 32-bit architecture if the order is ≥ 8.
func (h Hilbert4D) Pos(v []int) int {
	var d int
	N := 4
	for n := h.Order - 1; n >= 0; n-- {
		var e int
		for i := N - 1; i >= 0; i-- {
			v := v[i] >> n & 1
			e = e<<1 | (e^v)&1
		}

		d += e << (N * n)
		h.rot(false, h.Order, v, e)
	}
	return d
}

// Coord4D returns the spatial coordinates of pos.
func (h Hilbert4D) Coord4D(pos int) [4]int {
	N := 4
	var v [4]int
	for n := 0; n < h.Order; n++ {
		e := pos & (1<<N - 1)
		h.rot(true, n, v[:], e)

		for i, e := 0, e; i < N; i++ {
			v[i] += (e ^ e>>1) & 1 << n
			e >>= 1
		}
		pos >>= N
	}
	return v
}

// Coord returns the spatial coordinates of pos as a slice.
func (h Hilbert4D) Coord(pos int) []int {
	v := h.Coord4D(pos)
	return v[:]
}

type op interface{ do(int, []int) }

// invert I and J
type invert struct{ i, j int }

func (c invert) do(n int, v []int) { v[c.i], v[c.j] = v[c.i]^(1<<n-1), v[c.j]^(1<<n-1) }

// swap I and J
type swap struct{ i, j int }

func (c swap) do(n int, v []int) { v[c.i], v[c.j] = v[c.j], v[c.i] }

// swap and invert I and J
type flip struct{ i, j int }

func (c flip) do(n int, v []int) { v[c.i], v[c.j] = v[c.j]^(1<<n-1), v[c.i]^(1<<n-1) }

// do2 executes the given operations, optionally in reverse.
//
// Generic specialization reduces allocation (because it can eliminate interface
// value boxing) and improves performance
func do2[A, B op](reverse bool, n int, v []int, a A, b B) {
	if reverse {
		b.do(n, v)
		a.do(n, v)
	} else {
		a.do(n, v)
		b.do(n, v)
	}
}
