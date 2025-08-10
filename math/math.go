package math

import (
	"errors"
	"sort"
)

type NChooseR struct {
	N int
	R int
}

func (ncr NChooseR) Expand() (Quotient, error) {
	if ncr.R > ncr.N {
		return Quotient{}, errors.New("r cannot be greater than n")
	}
	n := Product{}
	for i := range ncr.R {
		n.NN = append(n.NN, ncr.N-i)
	}

	d := Product{}
	for i := range ncr.R {
		d.NN = append(d.NN, ncr.R-i)
	}

	return Quotient{A: n, B: d}, nil
}

type Product struct {
	NN []int
}

func (p Product) Mult(b Product) Product {
	return Product{NN: append(p.NN, b.NN...)}
}

type Quotient struct {
	A Product
	B Product
}

func (q Quotient) Mult(b Quotient) Quotient {
	return Quotient{
		A: q.A.Mult(b.A),
		B: q.B.Mult(b.B),
	}
}

func (q Quotient) Inverse() Quotient {
	return Quotient{
		A: q.B,
		B: q.A,
	}
}

func (q Quotient) Divide(b Quotient) Quotient {
	return q.Mult(b.Inverse())
}

func (q Quotient) Calculate() float64 {
	aa := q.A.NN
	bb := q.B.NN

	sort.Ints(aa)
	sort.Ints(bb)

	la := len(aa)
	lb := len(bb)

	l := la
	if lb > l {
		l = lb
	}

	cur := 1.0

	for i := range l {
		a := 1
		if i < la {
			a = aa[i]
		}

		b := 1
		if i < lb {
			b = bb[i]
		}

		cur *= (float64(a) / float64(b))
	}

	return cur
}
