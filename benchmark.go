package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"os"
	"runtime"
	"sync"
	"time"
)

var errInsufficientArgs = errors.New("insufficient arguments provided")

type config struct {
	name    string // Name of the cli application
	edgeUrl string // URL of the edge service
	vpsUrl  string // URL of the VPS service
	num     int64  // Number of requests to make to each service
}

type results struct {
	name  string
	times []time.Duration
	sum   int64
	p25   time.Duration
	p50   time.Duration
	p75   time.Duration
	p90   time.Duration
	p95   time.Duration
	p99   time.Duration
	max   time.Duration
	min   time.Duration
	avg   time.Duration
}

type benchmark map[string]*results

type pressureGauge struct {
	wg     sync.WaitGroup
	tokens chan struct{}
}

func newPressureGauge(limit int) *pressureGauge {
	ch := make(chan struct{}, limit)
	for range limit {
		ch <- struct{}{}
	}

	return &pressureGauge{
		tokens: ch,
	}
}

type sender func(int, string, string, benchmark)

func (pg *pressureGauge) execute(fn sender, i int, u string, s string, b benchmark) error {
	select {
	case <-pg.tokens:
		pg.wg.Add(1)
		go fn(i, u, s, b)
		return nil
	default:
		return errors.New("maximum requests in flight")
	}
}

var usageStr = `Compare the performance of a service deployed on a VPS vs on a serverless environment.

Usage: 
  %s [-h] [-n <value>] <edge-url> <vps-url>`

func parseArgs(w io.Writer, args []string) (*config, error) {
	var cfg config
	fs := flag.NewFlagSet("benchmark", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.Usage = func() {
		fmt.Fprintf(w, usageStr+"\n", fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Flags:")
		fs.PrintDefaults()
	}

	cfg.name = fs.Name()

	fs.Int64Var(&cfg.num, "n", 100, "Number of requests to make to each service.")

	if err := fs.Parse(args); err != nil {
		return &cfg, err
	}

	if fs.NArg() != 2 {
		return &cfg, errInsufficientArgs
	}

	cfg.edgeUrl = fs.Arg(0)
	cfg.vpsUrl = fs.Arg(1)

	return &cfg, nil
}

func validateConfig(config *config) error {
	if config.num < 1 {
		return errors.New("number of requests must be greater than 0")
	}

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	for _, url := range []string{config.edgeUrl, config.vpsUrl} {
		u, err := urlpkg.Parse(url)
		if err != nil {
			return err
		}

		if u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("invalid URL: %s", url)
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return fmt.Errorf("unsupported URL scheme: %s", u.Scheme)
		}

		_, err = client.Get(url)
		if err != nil {
			return err
		}
	}

	return nil
}

func send(i int, url string, server string, benchmark benchmark, pg *pressureGauge) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)
	if err != nil {
		return
	}
	resp.Body.Close()

	benchmark[server].times[i] = elapsed
	benchmark[server].sum += int64(elapsed)

	if elapsed > benchmark[server].max {
		benchmark[server].max = elapsed
	}

	if benchmark[server].min == 0 || elapsed < benchmark[server].min {
		benchmark[server].min = elapsed
	}

	pg.tokens <- struct{}{}
	pg.wg.Done()
}

func run(config *config, pg *pressureGauge, benchmark benchmark) error {
	var server string

	fmt.Printf("\nRunning benchmark...\n")
	benchmarkStartTime := time.Now()
	for i, url := range []string{config.edgeUrl, config.vpsUrl} {
		fmt.Printf("Making %d requests to %s...\n", config.num, url)
		switch i {
		case 0:
			server = "edge"
		case 1:
			server = "vps"
		}

		for i := 0; i < int(config.num); i++ {
			select {
			case <-pg.tokens:
				pg.wg.Add(1)
				go send(i, url, server, benchmark, pg)
			}
		}

		pg.wg.Wait()
		benchmark[server].avg = time.Duration(int64(benchmark[server].sum) / config.num)
	}

	fmt.Printf("\nBenchmark for %d reqs completed in %s\n", config.num, time.Since(benchmarkStartTime))
	fmt.Printf("Done.\n\n")

	fmt.Println("Results:")
	fmt.Printf("        avg          min          max\n")
	fmt.Printf("  Edge: %.3fms    %.3fms    %.3fms\n",
		benchmark["edge"].avg.Seconds()*1000, benchmark["edge"].min.Seconds()*1000, benchmark["edge"].max.Seconds()*1000)

	fmt.Printf("   VPS: %.3fms    %.3fms    %.3fms\n\n",
		benchmark["vps"].avg.Seconds()*1000, benchmark["vps"].min.Seconds()*1000, benchmark["vps"].max.Seconds()*1000)

	return nil
}

func main() {
	cfg, err := parseArgs(os.Stdout, os.Args[1:])
	if err != nil {
		if errors.Is(err, errInsufficientArgs) {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintf(os.Stderr, "Run '%s -h' for help.\n", cfg.name)
		}
		os.Exit(1)
	}

	if err := validateConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	benchmark := benchmark{
		"edge": &results{
			name:  "edge",
			times: make([]time.Duration, cfg.num),
		},
		"vps": &results{
			name:  "vps",
			times: make([]time.Duration, cfg.num),
		},
	}

	pg := newPressureGauge(runtime.GOMAXPROCS(0))

	if err := run(cfg, pg, benchmark); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
