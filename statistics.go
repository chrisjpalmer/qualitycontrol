package main

import (
	"iter"
	"maps"
	"math"
	"math/rand"
)

func sampleLotSTimes(n, k, r, s int) map[int]float64 {
	mcounts := make(map[int]int)
	for range s {
		m := sampleLot(n, k, r)
		mcounts[m]++
	}

	// calculate averages of different m(s)
	mdist := make(map[int]float64)
	for m, c := range mcounts {
		mdist[m] = float64(c) / float64(s)
	}

	return mdist
}

func sampleLot(n, k, r int) int {
	// create n parts
	nn := make([]bool, n)

	// mark k parts as broken
	for i := range k {
		nn[i] = true
	}

	// shuffle the parts around
	rand.Shuffle(n, func(i, j int) {
		nn[i], nn[j] = nn[j], nn[i]
	})

	// take a sample r
	rr := nn[:r]

	// count broken parts
	var m int
	for _, r := range rr {
		if r {
			m++
		}
	}

	return m
}

func klDivergence(pdist, mdist map[int]float64) float64 {

	mdistKeys := maps.Keys(mdist)
	pdistKeys := maps.Keys(pdist)

	min := minOfIters(mdistKeys, pdistKeys)
	max := maxOfIters(mdistKeys, pdistKeys)

	sum := 0.0
	for i := min; i < max; i++ {
		p := pdist[i]
		q := mdist[i]

		// Whenever P(x) is zero the contribution of the corresponding term
		// is interpreted as zero because lim of xlog(x) as x approaches infinity is 0
		if p == 0 {
			continue
		}

		// to avoid getting infinity, make q a very small number if 0
		if q == 0 {
			q = 0.00000001
		}

		sum += p * math.Log(p/q)
	}
	return sum
}

func minOfIters(iter1, iter2 iter.Seq[int]) int {
	min := 0
	for x := range iter1 {
		if x < min {
			min = x
		}
	}

	for x := range iter2 {
		if x < min {
			min = x
		}
	}

	return min
}

func maxOfIters(iter1, iter2 iter.Seq[int]) int {
	max := 0
	for x := range iter1 {
		if x > max {
			max = x
		}
	}

	for x := range iter2 {
		if x > max {
			max = x
		}
	}

	return max
}
