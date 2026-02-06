// Package convert provides character encoding conversion with parallel processing support.
// This file contains tests for the parallel processing functions.

package convert_test

import (
	"bytes"
	"testing"

	"github.com/bengarrett/retrotxtgo/convert"
	"github.com/nalgeon/be"
	"golang.org/x/text/encoding/charmap"
)

// Test ParallelConvert with different sizes
func TestParallelConvert(t *testing.T) {
	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Test with small data (should use regular conversion)
	smallData := []byte("Hello world!")
	runes, err := c.ParallelConvert(cp437, nil, smallData...)
	be.Err(t, err, nil)
	be.True(t, len(runes) > 0)

	// Test with medium data (should use chunked conversion)
	mediumData := bytes.Repeat([]byte("Hello world! "), 100) // ~1.2KB
	runes, err = c.ParallelConvert(cp437, nil, mediumData...)
	be.Err(t, err, nil)
	be.True(t, len(runes) > 0)

	// Test with large data (should use parallel conversion)
	largeData := bytes.Repeat([]byte("Hello world! This is a performance test. "), 1000) // ~36KB
	runes, err = c.ParallelConvert(cp437, nil, largeData...)
	be.Err(t, err, nil)
	be.True(t, len(runes) > 0)
}

// Test ChunkedConvert with different chunk sizes
func TestChunkedConvert(t *testing.T) {
	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Test with default chunk size
	data := bytes.Repeat([]byte("Hello world! "), 1000)
	runes, err := c.ChunkedConvert(cp437, nil, 0, data...)
	be.Err(t, err, nil)
	be.True(t, len(runes) > 0)

	// Test with custom chunk size
	runes, err = c.ChunkedConvert(cp437, nil, 4096, data...)
	be.Err(t, err, nil)
	be.True(t, len(runes) > 0)

	// Test with small chunk size
	runes, err = c.ChunkedConvert(cp437, nil, 1024, data...)
	be.Err(t, err, nil)
	be.True(t, len(runes) > 0)
}

// Test OptimalConvert automatic selection
func TestOptimalConvert(t *testing.T) {
	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Test with empty data
	runes, err := c.OptimalConvert(cp437, nil)
	be.Err(t, err, nil)
	be.Equal(t, len(runes), 0)

	// Test with small data (should use regular conversion)
	smallData := []byte("Hello")
	runes, err = c.OptimalConvert(cp437, nil, smallData...)
	be.Err(t, err, nil)
	be.True(t, len(runes) > 0)

	// Test with medium data (should use chunked conversion)
	mediumData := bytes.Repeat([]byte("Hello world! "), 500) // ~6KB
	runes, err = c.OptimalConvert(cp437, nil, mediumData...)
	be.Err(t, err, nil)
	be.True(t, len(runes) > 0)

	// Test with large data (should use parallel conversion)
	largeData := bytes.Repeat([]byte("Hello world! "), 10000) // ~120KB
	runes, err = c.OptimalConvert(cp437, nil, largeData...)
	be.Err(t, err, nil)
	be.True(t, len(runes) > 0)
}

// Benchmark parallel processing performance
func BenchmarkParallelProcessing(b *testing.B) {
	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Small data benchmark (should use regular conversion)
	smallData := bytes.Repeat([]byte("Hello world! "), 100) // ~1.2KB
	b.Run("SmallData", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = c.ParallelConvert(cp437, nil, smallData...)
		}
	})

	// Medium data benchmark (should use chunked conversion)
	mediumData := bytes.Repeat([]byte("Hello world! "), 1000) // ~12KB
	b.Run("MediumData", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = c.ParallelConvert(cp437, nil, mediumData...)
		}
	})

	// Large data benchmark (should use parallel conversion)
	largeData := bytes.Repeat([]byte("Hello world! "), 10000) // ~120KB
	b.Run("LargeData", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = c.ParallelConvert(cp437, nil, largeData...)
		}
	})

	// Very large data benchmark
	veryLargeData := bytes.Repeat([]byte("Hello world! "), 100000) // ~1.2MB
	b.Run("VeryLargeData", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = c.ParallelConvert(cp437, nil, veryLargeData...)
		}
	})
}

// Benchmark optimal conversion strategy
func BenchmarkOptimalConvert(b *testing.B) {
	cp437 := charmap.CodePage437
	c := convert.Convert{}

	// Test different data sizes
	sizes := []struct {
		name string
		size int
	}{
		{"Small", 1024},
		{"Medium", 8192},
		{"Large", 65536},
		{"VeryLarge", 1048576},
	}

	for _, size := range sizes {
		data := bytes.Repeat([]byte("Hello world! Testing performance. "), size.size)
		b.Run(size.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = c.OptimalConvert(cp437, nil, data...)
			}
		})
	}
}
