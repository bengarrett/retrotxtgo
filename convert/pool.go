package convert

import (
	"sync"
)

// Constants for memory pooling.
const (
	runeBufferInitialCapacity = 8192 // 8KB initial capacity for rune buffers
)

// runeBufferPool provides a pool of rune slices to reduce memory allocations
// and GC pressure in high-throughput scenarios.
//
//nolint:gochecknoglobals // Intentional global pool for performance optimization
var runeBufferPool = &sync.Pool{
	New: func() any {
		return make([]rune, 0, runeBufferInitialCapacity) // runeBufferInitialCapacity initial capacity
	},
}

// getRuneBuffer acquires a rune buffer from the pool or creates a new one.
func getRuneBuffer() []rune {
	if buf, ok := runeBufferPool.Get().([]rune); ok {
		return buf
	}
	// If type assertion fails, return a new buffer
	return make([]rune, 0, runeBufferInitialCapacity)
}
