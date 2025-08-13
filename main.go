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

	s, err := strconv.Atoi(q.Get("s"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	line := setupChart()

	// calculate probabilities
	mm, pp, min, max, err := pMMDefective(n, k, r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	mmdata := []opts.LineData{}
	fmt.Printf("\n----------------------------------------------------\nn: %d, k: %d, r: %d\n", n, k, r)
	for i, p := range pp {
		fmt.Printf("m: %d => %.2f%%\n", mm[i], p*100)
		mmdata = append(mmdata, opts.LineData{Value: p})
	}

	line.SetXAxis(mm).
		AddSeries("P(m defective in r samples)", mmdata).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}))

	// do experiment and graph the results
	statMM := sampleLotSTimes(n, k, r, s)

	if !mInBounds(statMM, min, max) {
		w.WriteHeader(400)
		w.Write([]byte("m not in bounds! What have you done Chris?"))
		return
	}

	fmt.Printf("\n----------------------------------------------------\nstats\n\nn: %d, k: %d, r: %d\n", n, k, r)

	statMMData := []opts.LineData{}

	for _, m := range mm {
		avg := statMM[m]
		fmt.Printf("m: %d => %.2f%%\n", m, avg*100)
		statMMData = append(statMMData, opts.LineData{Value: avg})
	}

	line.AddSeries("Avg(m defective in r samples)", statMMData).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(false)}))

	line.Render(w)

	avgMse := avgMse(n, k, r, s, mm, pp)

	fmt.Printf("avg mse: %.10f%%\n", avgMse)

	writeAvgMSE(w, avgMse)
}

func writeAvgMSE(w io.Writer, avgMse float64) {
	s := fmt.Sprintf(`<br><label>avg mse: %0.10f%%</label>`, avgMse)
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
