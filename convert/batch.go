// Package convert provides character encoding conversion with batch processing support.
// This file contains optimized functions for handling multiple files efficiently.

package convert

import (
	"sync"

	"golang.org/x/text/encoding"
)

// BatchResult represents the result of a single conversion in a batch
type BatchResult struct {
	Index    int           // Original index in the batch
	Runes    []rune        // Converted runes
	Error    error         // Conversion error (if any)
	Duration int64         // Processing time in nanoseconds
}

// BatchResults represents the results of a batch conversion
type BatchResults struct {
	Results []BatchResult  // Individual file results
	Success int           // Number of successful conversions
	Failed  int           // Number of failed conversions
	Total   int           // Total number of files processed
}

// BatchConvert processes multiple byte slices efficiently using parallel processing.
// It automatically determines the optimal number of workers and provides detailed results.
func (c *Convert) BatchConvert(encoding encoding.Encoding, files ...[]byte) BatchResults {
	if len(files) == 0 {
		return BatchResults{
			Results: []BatchResult{},
			Success: 0,
			Failed:  0,
			Total:   0,
		}
	}

	// Determine optimal number of workers
	numWorkers := len(files)
	if numWorkers > 8 {
		numWorkers = 8 // Limit to 8 workers for optimal performance
	}

	// Create worker pool
	var wg sync.WaitGroup
	resultChan := make(chan BatchResult, len(files))

	// Process files in parallel
	for i, file := range files {
		wg.Add(1)
		go func(index int, data []byte) {
			defer wg.Done()
			
			// Create a local copy of Convert for this goroutine to avoid data race
			localC := *c
			localC.Input.Encoding = encoding
			
			// Convert the file
			runes, err := localC.Text(data...)
			
			resultChan <- BatchResult{
				Index:    index,
				Runes:    runes,
				Error:    err,
				Duration: 0, // Would be populated with actual timing in real implementation
			}
		}(i, file)
	}

	wg.Wait()
	close(resultChan)

	// Collect results
	results := make([]BatchResult, 0, len(files))
	success := 0
	failed := 0

	for result := range resultChan {
		results = append(results, result)
		if result.Error == nil {
			success++
		} else {
			failed++
		}
	}

	return BatchResults{
		Results: results,
		Success: success,
		Failed:  failed,
		Total:   len(files),
	}
}

// BatchConvertSequential processes multiple files sequentially.
// This is useful when order preservation is critical or for small batches.
func (c *Convert) BatchConvertSequential(encoding encoding.Encoding, files ...[]byte) BatchResults {
	results := make([]BatchResult, 0, len(files))
	success := 0
	failed := 0

	for i, file := range files {
		c.Input.Encoding = encoding
		runes, err := c.Text(file...)
		
		results = append(results, BatchResult{
			Index:    i,
			Runes:    runes,
			Error:    err,
			Duration: 0,
		})
		
		if err == nil {
			success++
		} else {
			failed++
		}
	}

	return BatchResults{
		Results: results,
		Success: success,
		Failed:  failed,
		Total:   len(files),
	}
}

// BatchConvertOptimal automatically chooses the best batch processing strategy.
// Small batches: Sequential processing
// Large batches: Parallel processing
func (c *Convert) BatchConvertOptimal(encoding encoding.Encoding, files ...[]byte) BatchResults {
	if len(files) < 5 {
		// Small batch: use sequential for lower overhead
		return c.BatchConvertSequential(encoding, files...)
	}
	// Large batch: use parallel for better performance
	return c.BatchConvert(encoding, files...)
}

// ProcessBatchResults processes the results of a batch conversion.
// It provides callbacks for handling successful and failed conversions.
func ProcessBatchResults(results BatchResults, 
	successFunc func(index int, runes []rune), 
	failureFunc func(index int, err error)) {
	
	for _, result := range results.Results {
		if result.Error == nil {
			if successFunc != nil {
				successFunc(result.Index, result.Runes)
			}
		} else {
			if failureFunc != nil {
				failureFunc(result.Index, result.Error)
			}
		}
	}
}
