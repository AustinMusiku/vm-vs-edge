package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"os"
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

type benchmark struct {
	edge *results
	vps  *results
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

func run(config *config) error {
	edgeResults := &results{
		name:  "edge",
		times: make([]time.Duration, config.num),
	}

	vpsResults := &results{
		name:  "vps",
		times: make([]time.Duration, config.num),
	}

	benchmark := benchmark{
		edge: edgeResults,
		vps:  vpsResults,
	}

	fmt.Printf("\nRunning benchmark...\n")
	for _, url := range []string{config.edgeUrl, config.vpsUrl} {
		fmt.Printf("Making %d requests to %s...\n", config.num, url)

		client := http.Client{
			Timeout: 3 * time.Second,
		}

		for i := 0; i < int(config.num); i++ {
			start := time.Now()
			resp, err := client.Get(url)
			stop := time.Since(start)
			if err != nil {
				return err
			}

			resp.Body.Close()

			if url == config.edgeUrl {
				edgeResults.times[i] = stop
				edgeResults.sum += int64(stop)
			} else {
				vpsResults.times[i] = stop
				vpsResults.sum += int64(stop)
			}
		}

		if url == config.edgeUrl {
			edgeResults.avg = time.Duration(int64(edgeResults.sum) / config.num)
		} else {
			vpsResults.avg = time.Duration(int64(vpsResults.sum) / config.num)

		}
	}

	fmt.Printf("  Done.\n\n")

	fmt.Printf("Edge avg: %s\n", benchmark.edge.avg)
	fmt.Printf(" VPS avg: %s\n", benchmark.vps.avg)

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

	fmt.Println("\nConfiguration:")
	fmt.Printf("  Edge URL: %s\n", cfg.edgeUrl)
	fmt.Printf("   VPS URL: %s\n", cfg.vpsUrl)
	fmt.Printf("  Requests: %d\n", cfg.num)

	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
