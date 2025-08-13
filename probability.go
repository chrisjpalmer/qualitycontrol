package main

import (
	"errors"

	qmath "github.com/chrisjpalmer/qualitycontrol/math"
)

func pMMDefective(n, k, r int) (mm []int, pp []float64, min int, max int, err error) {
	// calculate distribution
	mm = []int{}
	pp = []float64{}

	max = k
	if r < max {
		max = r
	}

	for i := min; i < max+1; i++ {
		mm = append(mm, i)

		p, err := pMDefective(n, k, r, i)
		if err != nil {
			return nil, nil, 0, 0, err
		}

		if p > 1 {
			return nil, nil, 0, 0, errors.New("floating point math error")
		}

		pp = append(pp, p)
	}

	return mm, pp, min, max, nil
}

func pMDefective(n, k, r, m int) (float64, error) {
	chooseGood, err := qmath.NChooseR{N: (n - k), R: (r - m)}.Expand()
	if err != nil {
		return 0, err
	}

	chooseBad, err := qmath.NChooseR{N: k, R: m}.Expand()
	if err != nil {
		return 0, err
	}

	totalSamples, err := qmath.NChooseR{N: n, R: r}.Expand()
	if err != nil {
		return 0, err
	}

	pq := chooseGood.Mult(chooseBad).Divide(totalSamples)

	p := pq.Calculate()

	return p, nil
}
