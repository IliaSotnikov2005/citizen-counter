package counter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
)

type ParallelCounter struct {
	Workers      int
	BufferSizeMB int
}

type filePart struct {
	offset int64
	size   int64
}

type partResult struct {
	index  int
	counts map[string]int
	err    error
}

func NewParallelCounter(workers, bufferMB int) *ParallelCounter {
	return &ParallelCounter{
		Workers:      workers,
		BufferSizeMB: bufferMB,
	}
}

func splitFile(filePath string, numParts int) ([]filePart, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	fileSize := info.Size()
	if fileSize == 0 {
		return nil, nil
	}

	partSize := max(fileSize/int64(numParts), 1024)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	parts := make([]filePart, 0, numParts)
	start := int64(0)

	for start < fileSize {
		end := start + partSize
		if end >= fileSize {
			end = fileSize
		} else {
			newEnd, err := findLineEnd(file, end, fileSize)
			if err != nil {
				return nil, err
			}

			end = newEnd
		}

		parts = append(parts, filePart{
			offset: start,
			size:   end - start,
		})
		start = end
	}

	return parts, nil
}

func findLineEnd(file *os.File, pos, fileSize int64) (int64, error) {
	if pos >= fileSize {
		return fileSize, nil
	}

	_, err := file.Seek(pos, 0)
	if err != nil {
		return pos, err
	}

	buf := make([]byte, 4096)
	currentPos := pos

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return currentPos, err
		}
		if n == 0 {
			break
		}

		for i := range n {
			currentPos++
			if buf[i] == '\n' {
				return currentPos, nil
			}
		}
		if err == io.EOF {
			break
		}
	}

	return fileSize, nil
}

func (c *ParallelCounter) Count(filePath string) (*Result, error) {
	if c.Workers <= 0 {
		c.Workers = runtime.NumCPU()
	}

	parts, err := splitFile(filePath, c.Workers)
	if err != nil {
		return nil, fmt.Errorf("failed to split file: %w", err)
	}
	if len(parts) == 0 {
		return &Result{Counts: make(map[string]int)}, nil
	}

	resultsCh := make(chan *partResult, len(parts))
	var wg sync.WaitGroup

	for i, part := range parts {
		wg.Add(1)
		go func(partIndex int, p filePart) {
			defer wg.Done()
			counts, err := processPart(filePath, p, c.BufferSizeMB)
			resultsCh <- &partResult{
				index:  partIndex,
				counts: counts,
				err:    err,
			}
		}(i, part)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	finalCounts := make(map[string]int, estimateCapacity(parts[0].size*int64(len(parts))))
	for result := range resultsCh {
		if result.err != nil {
			return nil, fmt.Errorf("part %d failed: %w", result.index, result.err)
		}
		for name, count := range result.counts {
			finalCounts[name] += count
		}
	}

	return &Result{Counts: finalCounts}, nil
}

func processPart(filePath string, part filePart, bufferMB int) (map[string]int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Seek(part.offset, 0)
	if err != nil {
		return nil, err
	}

	limitedReader := io.LimitReader(file, part.size)
	bufSize := max(bufferMB*1024*1024, 64*1024)
	reader := bufio.NewReaderSize(limitedReader, bufSize)

	counts := make(map[string]int, 1000)

	for {
		line, err := reader.ReadString('\n')
		name := strings.TrimSpace(line)
		if name != "" {
			counts[name]++
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return counts, nil
}
