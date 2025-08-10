package math_test

import (
	"testing"

	"github.com/chrisjpalmer/qualitycontrol/math"
	"github.com/stretchr/testify/assert"
)

func TestProductMult(t *testing.T) {
	assert.Equal(t, math.Product{NN: []int{2}}.Mult(math.Product{NN: []int{3}}).NN, []int{2, 3})
}

func TestQuotient(t *testing.T) {
	assert.Equal(t, math.Quotient{A: math.Product{NN: []int{5, 4}}, B: math.Product{NN: []int{5, 4}}}.Calculate(), float64(1))
	assert.Equal(t, math.Quotient{A: math.Product{NN: []int{5, 4}}, B: math.Product{NN: []int{4, 5}}}.Calculate(), float64(1))
	assert.Equal(t, math.Quotient{A: math.Product{NN: []int{5, 4}}, B: math.Product{NN: []int{1}}}.Calculate(), float64(20))
}

func TestQuotientMult(t *testing.T) {
	a := math.Quotient{A: math.Product{NN: []int{2}}, B: math.Product{NN: []int{3}}}
	b := math.Quotient{A: math.Product{NN: []int{3}}, B: math.Product{NN: []int{2}}}
	assert.Equal(t, a.Mult(b).Calculate(), 1.0)
}

func TestQuotientDivide(t *testing.T) {
	a := math.Quotient{A: math.Product{NN: []int{2}}, B: math.Product{NN: []int{3}}}
	b := math.Quotient{A: math.Product{NN: []int{2}}, B: math.Product{NN: []int{3}}}
	assert.Equal(t, a.Divide(b).Calculate(), 1.0)
}

func TestNChooseR(t *testing.T) {
	res, err := math.NChooseR{N: 5, R: 3}.Expand()
	assert.NoError(t, err)
	assert.Equal(t, math.Quotient{A: math.Product{NN: []int{5, 4, 3}}, B: math.Product{NN: []int{3, 2, 1}}}, res)

	res, err = math.NChooseR{N: 5, R: 5}.Expand()
	assert.NoError(t, err)
	assert.Equal(t, math.Quotient{A: math.Product{NN: []int{5, 4, 3, 2, 1}}, B: math.Product{NN: []int{5, 4, 3, 2, 1}}}, res)

	res, err = math.NChooseR{N: 5, R: 1}.Expand()
	assert.NoError(t, err)
	assert.Equal(t, math.Quotient{A: math.Product{NN: []int{5}}, B: math.Product{NN: []int{1}}}, res)
}
