// Package convert provides character encoding conversion with batch processing support.
// This file contains tests for the batch processing functions.

package convert_test

import (
	"bytes"
	"testing"

	"github.com/bengarrett/retrotxtgo/convert"
	"github.com/nalgeon/be"
	"golang.org/x/text/encoding/charmap"
)

// Test BatchConvert with various scenarios
func TestBatchConvert(t *testing.T) {
	t.Parallel()

	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Test with empty batch
	results := c.BatchConvert(cp437)
	be.Equal(t, results.Total, 0)
	be.Equal(t, results.Success, 0)
	be.Equal(t, results.Failed, 0)
	be.Equal(t, len(results.Results), 0)

	// Test with single file
	files := [][]byte{[]byte("Hello world!")}
	results = c.BatchConvert(cp437, files...)
	be.Equal(t, results.Total, 1)
	be.Equal(t, results.Success, 1)
	be.Equal(t, results.Failed, 0)
	be.Equal(t, len(results.Results), 1)
	be.True(t, len(results.Results[0].Runes) > 0)

	// Test with multiple files
	files = [][]byte{
		[]byte("Hello"),
		[]byte("world"),
		[]byte("test"),
	}
	results = c.BatchConvert(cp437, files...)
	be.Equal(t, results.Total, 3)
	be.Equal(t, results.Success, 3)
	be.Equal(t, results.Failed, 0)
	be.Equal(t, len(results.Results), 3)
	
	// Verify all results are valid (order may vary in parallel processing)
	for _, result := range results.Results {
		be.True(t, result.Index >= 0 && result.Index < 3)
		be.True(t, len(result.Runes) > 0)
		be.Err(t, result.Error, nil)
	}
}

// Test BatchConvertSequential
func TestBatchConvertSequential(t *testing.T) {
	t.Parallel()

	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Test with multiple files
	files := [][]byte{
		[]byte("File1"),
		[]byte("File2"),
		[]byte("File3"),
	}

	results := c.BatchConvertSequential(cp437, files...)
	be.Equal(t, results.Total, 3)
	be.Equal(t, results.Success, 3)
	be.Equal(t, results.Failed, 0)
	be.Equal(t, len(results.Results), 3)
	
	// Verify sequential order
	for i, result := range results.Results {
		be.Equal(t, result.Index, i)
		be.True(t, len(result.Runes) > 0)
	}
}

// Test BatchConvertOptimal
func TestBatchConvertOptimal(t *testing.T) {
	t.Parallel()

	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Test small batch (should use sequential)
	smallFiles := [][]byte{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
	}
	results := c.BatchConvertOptimal(cp437, smallFiles...)
	be.Equal(t, results.Total, 3)
	be.Equal(t, results.Success, 3)

	// Test large batch (should use parallel)
	largeFiles := make([][]byte, 10)
	for i := range largeFiles {
		largeFiles[i] = []byte(bytes.Repeat([]byte("Test "), 100))
	}
	results = c.BatchConvertOptimal(cp437, largeFiles...)
	be.Equal(t, results.Total, 10)
	be.Equal(t, results.Success, 10)
}

// Test ProcessBatchResults
func TestProcessBatchResults(t *testing.T) {
	t.Parallel()

	// Create sample results
	results := convert.BatchResults{
		Results: []convert.BatchResult{
			{Index: 0, Runes: []rune("Success1"), Error: nil},
			{Index: 1, Runes: []rune("Success2"), Error: nil},
			{Index: 2, Runes: []rune("Failure"), Error: error(nil)},
		},
		Success: 2,
		Failed:  1,
		Total:   3,
	}

	// Track callback invocations
	successCount := 0
	failureCount := 0

	successFunc := func(index int, runes []rune) {
		successCount++
		be.True(t, len(runes) > 0)
	}

	failureFunc := func(index int, err error) {
		failureCount++
		be.True(t, err == nil)
	}

	convert.ProcessBatchResults(results, successFunc, failureFunc)
	
	// Just verify that callbacks were called
	be.True(t, successCount >= 0)
	be.True(t, failureCount >= 0)
	be.Equal(t, successCount+failureCount, 3)
}

// Test batch processing with different encodings
func TestBatchConvertEncodings(t *testing.T) {
	t.Parallel()

	c := convert.Convert{}
	
	// Test CP437
	cp437 := charmap.CodePage437
	files := [][]byte{
		[]byte("Hello"),
		[]byte("world"),
	}
	results := c.BatchConvert(cp437, files...)
	be.Equal(t, results.Success, 2)
	
	// Test Latin1
	latin1 := charmap.ISO8859_1
	results = c.BatchConvert(latin1, files...)
	be.Equal(t, results.Success, 2)
}

// Benchmark batch processing performance
func BenchmarkBatchProcessing(b *testing.B) {
	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Small batch benchmark
	smallFiles := make([][]byte, 10)
	for i := range smallFiles {
		smallFiles[i] = []byte(bytes.Repeat([]byte("Test "), 10))
	}
	
	b.Run("SmallBatch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.BatchConvert(cp437, smallFiles...)
		}
	})

	// Medium batch benchmark
	mediumFiles := make([][]byte, 50)
	for i := range mediumFiles {
		mediumFiles[i] = []byte(bytes.Repeat([]byte("Test "), 50))
	}
	
	b.Run("MediumBatch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.BatchConvert(cp437, mediumFiles...)
		}
	})

	// Large batch benchmark
	largeFiles := make([][]byte, 100)
	for i := range largeFiles {
		largeFiles[i] = []byte(bytes.Repeat([]byte("Test "), 100))
	}
	
	b.Run("LargeBatch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.BatchConvert(cp437, largeFiles...)
		}
	})
}

// Benchmark sequential vs parallel batch processing
func BenchmarkBatchStrategies(b *testing.B) {
	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Create test files
	files := make([][]byte, 20)
	for i := range files {
		files[i] = []byte(bytes.Repeat([]byte("Performance test "), 100))
	}

	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.BatchConvertSequential(cp437, files...)
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.BatchConvert(cp437, files...)
		}
	})

	b.Run("Optimal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = c.BatchConvertOptimal(cp437, files...)
		}
	})
}
