package main

import (
	_ "net/http/pprof"

	"flag"
	"fmt"
	"os"

	"github.com/IliaSotnikov2005/citizen-counter/internal/counter"
	"github.com/IliaSotnikov2005/citizen-counter/internal/output"
)

const (
	sequentialMode = "sequential"
	parallelMode   = "parallel"
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

	var citizenCounter counter.Counter
	switch *mode {
	case sequentialMode:
		citizenCounter = counter.NewSequentialCounter(*bufSize)
	case parallelMode:
		citizenCounter = counter.NewParallelCounter(*workers, *bufSize)
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s. Use 'sequential' or 'parallel'.\n", *mode)
		os.Exit(1)
	}

	result, err := citizenCounter.Count(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	output.Print(result.Counts, *topN, *sort)
}
