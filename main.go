package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	_ "embed"

	"github.com/chrisjpalmer/qualitycontrol/math"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

//go:embed index.html
var indexPageBytes []byte

func indexPage(w http.ResponseWriter, rq *http.Request) {
	w.Write(indexPageBytes)
}

func chart(w http.ResponseWriter, rq *http.Request) {

	q := rq.URL.Query()

	n, err := strconv.Atoi(q.Get("n"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	k, err := strconv.Atoi(q.Get("k"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	r, err := strconv.Atoi(q.Get("r"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	// create a new line instance
	line := charts.NewLine()

	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeChalk}),
		charts.WithTitleOpts(opts.Title{
			Title: "P(m defective in r samples)",
		}),
		charts.WithAnimation(false),
	)

	mm, pp, err := calcProbabilityDistributionM(n, k, r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	// make data points
	dp := []opts.LineData{}
	for i, p := range pp {
		fmt.Printf("n: %d, k: %d, r: %d, m: %d => %.2f%%\n", n, k, r, mm[i], p*100)
		dp = append(dp, opts.LineData{Value: p})
	}

	// Put data into instance
	line.SetXAxis(mm).
		AddSeries("P(m defective in samples size of R)", dp).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}))
	line.Render(w)
}

func main() {
	http.DefaultServeMux.HandleFunc("/chart", chart)
	http.DefaultServeMux.HandleFunc("/", indexPage)
	http.ListenAndServe(":8081", nil)
}

func calcProbabilityDistributionM(n, k, r int) ([]int, []float64, error) {
	// calculate distribution
	mm := []int{}
	pp := []float64{}

	max := k
	if r < max {
		max = r
	}

	for i := range max + 1 {
		mm = append(mm, i)

		p, err := calcProbabilityMInSampleSizeR(n, k, r, i)
		if err != nil {
			return nil, nil, err
		}

		if p > 1 {
			return nil, nil, errors.New("floating point math error")
		}

		pp = append(pp, p)
	}

	return mm, pp, nil
}

func calcProbabilityMInSampleSizeR(n, k, r, m int) (float64, error) {
	chooseGood, err := math.NChooseR{N: (n - k), R: (r - m)}.Expand()
	if err != nil {
		return 0, err
	}

	chooseBad, err := math.NChooseR{N: k, R: m}.Expand()
	if err != nil {
		return 0, err
	}

	totalSamples, err := math.NChooseR{N: n, R: r}.Expand()
	if err != nil {
		return 0, err
	}

	pq := chooseGood.Mult(chooseBad).Divide(totalSamples)

	p := pq.Calculate()

	return p, nil
}
