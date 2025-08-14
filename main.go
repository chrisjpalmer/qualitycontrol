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
	n, k, r, n2, k2, r2, s int
	compare                bool
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

	qp.n2, err = strconv.Atoi(q.Get("n2"))
	if err != nil {
		return QueryParams{}, err
	}

	qp.k2, err = strconv.Atoi(q.Get("k2"))
	if err != nil {
		return QueryParams{}, err
	}

	qp.r2, err = strconv.Atoi(q.Get("r2"))
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
	mm, pp, min, max, err := pMMDefective(qp.n, qp.k, qp.r)
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
		statMM := sampleLotSTimes(qp.n2, qp.k2, qp.r2, qp.s)

		if !mInBounds(statMM, min, max) {
			w.WriteHeader(400)
			w.Write([]byte("m not in bounds! What have you done Chris?"))
			return
		}

		fmt.Printf("\n----------------------------------------------------\nstats\n\nn: %d, k: %d, r: %d\n", qp.n2, qp.k2, qp.r2)

		statMMData := []opts.LineData{}

		for _, m := range mm {
			avg := statMM[m]
			fmt.Printf("m: %d => %.2f%%\n", m, avg*100)
			statMMData = append(statMMData, opts.LineData{Value: avg})
		}

		line.AddSeries("Avg(m defective in r samples)", statMMData).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}))
	}

	line.Render(w)

	if qp.compare {
		mse := mse(qp.n2, qp.k2, qp.r2, qp.s, mm, pp)

		fmt.Printf("mse: %.10f%%\n", mse)

		writeMSE(w, mse)
	}
}

func writeMSE(w io.Writer, mse float64) {
	s := fmt.Sprintf(`<br><label>mse: %0.10f%%</label>`, mse)
	w.Write([]byte(s))
}

func setupChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeChalk}),
		charts.WithTitleOpts(opts.Title{
			Title: "P(m defective in r samples)",
		}),
		charts.WithAnimation(false),
		charts.WithYAxisOpts(opts.YAxis{Min: 0, Max: 1}),
	)

	return line
}
