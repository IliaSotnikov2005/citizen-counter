package main

import (
	_ "net/http/pprof"

	"flag"
	"fmt"
	"os"

	"github.com/IliaSotnikov2005/citizen-counter/internal/counter"
	"github.com/IliaSotnikov2005/citizen-counter/internal/output"
)

var (
	mode    = flag.String("mode", "sequential", "processing mode: sequential or parallel")
	workers = flag.Int("workers", 0, "number of workers for parallel mode (0 = CPU cores)")
)

func main() {
	topN := flag.Int("top", 20, "show first N (0 = all)")
	bufSize := flag.Int("buf", 64, "buffer size in MB")
	sort := flag.Bool("sort", true, "sort by frequency")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	filePath := flag.Arg(0)

	var result *counter.Result
	var err error

	if *mode == "parallel" {
		fmt.Fprintf(os.Stderr, "Running in PARALLEL mode with %d workers\n", *workers)
		config := counter.ParallelConfig{
			Workers: *workers,
			BufSize: *bufSize,
		}

		result, err = counter.CountParallel(filePath, config)
	} else {
		fmt.Fprintf(os.Stderr, "Running in SEQUENTIAL mode\n")
		result, err = counter.Count(filePath, *bufSize)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	output.Print(result.Counts, *topN, *sort)
}
