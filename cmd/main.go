package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/IliaSotnikov2005/citizen-counter/internal/counter"
	"github.com/IliaSotnikov2005/citizen-counter/internal/output"
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

	result, err := counter.Count(filePath, *bufSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	output.Print(result.Counts, *topN, *sort)
}
