package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	_ "embed"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func main() {
	http.DefaultServeMux.HandleFunc("/chart", chart)
	http.DefaultServeMux.HandleFunc("/", indexPage)

	done := make(chan struct{})
	go func() {
		defer close(done)
		http.ListenAndServe(":8081", nil)
	}()

	fmt.Println("server is running! listening on port 8081")

	<-done
}

//go:embed index.html
var indexPageBytes []byte

func indexPage(w http.ResponseWriter, rq *http.Request) {
	w.Write(indexPageBytes)
}

type QueryParams struct {
	n, k, r, k2, s int
	compare        bool
}

func parseQuery(rq *http.Request) (QueryParams, error) {
	q := rq.URL.Query()

	var err error
	var qp QueryParams

	qp.n, err = strconv.Atoi(q.Get("n"))
	if err != nil {
		return QueryParams{}, err
	}

	qp.k, err = strconv.Atoi(q.Get("k"))
	if err != nil {
		return QueryParams{}, err
	}

	qp.r, err = strconv.Atoi(q.Get("r"))
	if err != nil {
		return QueryParams{}, err
	}

	qp.k2, err = strconv.Atoi(q.Get("k2"))
	if err != nil {
		return QueryParams{}, err
	}

	qp.s, err = strconv.Atoi(q.Get("s"))
	if err != nil {
		return QueryParams{}, err
	}

	if q.Get("compare") == "true" {
		qp.compare = true
	}

	return qp, nil
}

func chart(w http.ResponseWriter, rq *http.Request) {

	line := setupChart()

	qp, err := parseQuery(rq)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	// calculate probabilities
	mm, pp, pdist, err := pMMDefective(qp.n, qp.k, qp.r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	mmdata := []opts.LineData{}
	fmt.Printf("\n----------------------------------------------------\nn: %d, k: %d, r: %d\n", qp.n, qp.k, qp.r)
	for i, p := range pp {
		fmt.Printf("m: %d => %.2f%%\n", mm[i], p*100)
		mmdata = append(mmdata, opts.LineData{Value: p})
	}

	line.SetXAxis(mm).
		AddSeries("P(m defective in r samples)", mmdata).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}))

	if qp.compare {
		// do experiment and graph the results
		statMM := sampleLotSTimes(qp.n, qp.k2, qp.r, qp.s)

		fmt.Printf("\n----------------------------------------------------\nstats\n\nn: %d, k: %d, r: %d\n", qp.n, qp.k2, qp.r)

		statMMData := []opts.LineData{}

		for _, m := range mm {
			avg := statMM[m]
			fmt.Printf("m: %d => %.2f%%\n", m, avg*100)
			statMMData = append(statMMData, opts.LineData{Value: avg})
		}

		line.AddSeries(fmt.Sprintf("Avg(m defective in r samples over %d experiments)", qp.s), statMMData).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}))
	}

	line.Render(w)

	if qp.compare {
		mdist := sampleLotSTimes(qp.n, qp.k2, qp.r, qp.s)

		kl := klDivergence(pdist, mdist)

		fmt.Printf("kl: %.4f\n", kl)

		writeKLDivergence(w, kl)
	}
}

func writeKLDivergence(w io.Writer, kl float64) {
	s := fmt.Sprintf(`<br><label>kl divergence: %0.4f</label>`, kl)
	w.Write([]byte(s))
}

func setupChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeChalk}),
		charts.WithAnimation(false),
		charts.WithYAxisOpts(opts.YAxis{Min: 0, Max: 1}),
	)

	return line
}
