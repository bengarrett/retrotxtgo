// Package convert provides character encoding conversion with parallel processing support.
// This file contains optimized functions for handling large file conversions.

package convert

import (
	"runtime"
	"sync"

	"golang.org/x/text/encoding"
)

// ParallelConvert processes large byte slices using parallel processing.
// It automatically determines the optimal number of workers based on the
// input size and available CPU cores.
func (c *Convert) ParallelConvert(in, out encoding.Encoding, b ...byte) ([]rune, error) {
	// For small inputs, use regular conversion
	if len(b) < 8192 { // 8KB threshold
		c.Input.Encoding = in
		return c.Text(b...)
	}

	// Determine optimal chunk size and number of workers
	numWorkers := runtime.NumCPU()
	if numWorkers > 4 {
		numWorkers = 4 // Limit to 4 workers for optimal performance
	}

	chunkSize := (len(b) + numWorkers - 1) / numWorkers

	// Create worker pool
	var wg sync.WaitGroup
	results := make([][]rune, numWorkers)
	errChan := make(chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			start := workerID * chunkSize
			end := start + chunkSize
			if workerID == numWorkers-1 {
				end = len(b)
			}

			if start >= end {
				return
			}

			// Create a local copy of Convert for this goroutine to avoid data race
			localC := *c
			localC.Input.Encoding = in
			runes, err := localC.Text(b[start:end]...)
			if err != nil {
				errChan <- err
				return
			}
			results[workerID] = runes
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	if len(errChan) > 0 {
		return nil, <-errChan
	}

	// Combine results
	totalRunes := 0
	for _, r := range results {
		totalRunes += len(r)
	}

	combined := make([]rune, 0, totalRunes)
	for _, r := range results {
		combined = append(combined, r...)
	}

	return combined, nil
}

// ChunkedConvert processes data in fixed-size chunks for memory efficiency.
// This is useful for very large files where memory usage is a concern.
func (c *Convert) ChunkedConvert(in, out encoding.Encoding, chunkSize int, b ...byte) ([]rune, error) {
	if chunkSize <= 0 {
		chunkSize = 8192 // Default 8KB chunks
	}

	var result []rune
	for i := 0; i < len(b); i += chunkSize {
		end := i + chunkSize
		if end > len(b) {
			end = len(b)
		}

		c.Input.Encoding = in
		runes, err := c.Text(b[i:end]...)
		if err != nil {
			return nil, err
		}
		result = append(result, runes...)
	}

	return result, nil
}

// OptimalConvert automatically chooses the best conversion method based on input size.
// Small inputs: Regular conversion
// Medium inputs: Chunked conversion
// Large inputs: Parallel conversion
func (c *Convert) OptimalConvert(in, out encoding.Encoding, b ...byte) ([]rune, error) {
	data := b
	if len(data) == 0 {
		return nil, nil
	}

	// Small data: use regular conversion
	if len(data) < 8192 { // 8KB
		c.Input.Encoding = in
		return c.Text(data...)
	}

	// Medium data: use chunked conversion
	if len(data) < 1024*1024 { // 1MB
		return c.ChunkedConvert(in, out, 32768, data...) // 32KB chunks
	}

	// Large data: use parallel conversion
	return c.ParallelConvert(in, out, data...)
}
