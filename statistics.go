package main

import (
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
	mavg := make(map[int]float64)
	for m, c := range mcounts {
		mavg[m] = float64(c) / float64(s)
	}

	return mavg
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

func mInBounds(statMM map[int]float64, min int, max int) bool {
	for m := range statMM {
		if m < min || m > max {
			return false
		}
	}

	return true
}

func mse(n, k, r, s int, mm []int, pp []float64) float64 {
	statMM := sampleLotSTimes(n, k, r, s)

	sum := 0.0
	for _, m := range mm {
		s := statMM[m]
		p := pp[m]
		sum += math.Pow(s-p, 2)
	}
	return (sum / float64(len(mm))) * 100
}
