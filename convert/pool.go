// Package convert provides character encoding conversion with memory pooling support.
// This file contains memory pooling optimizations to reduce GC pressure in high-throughput scenarios.

package convert

import (
	"sync"
)

// runeBufferPool provides a pool of rune slices to reduce memory allocations
// and GC pressure in high-throughput scenarios.
var runeBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]rune, 0, 8192) // 8KB initial capacity
	},
}

// getRuneBuffer acquires a rune buffer from the pool or creates a new one.
func getRuneBuffer() []rune {
	return runeBufferPool.Get().([]rune)
}

// putRuneBuffer returns a rune buffer to the pool for reuse.
func putRuneBuffer(buf []rune) {
	// Clear the buffer before returning to pool
	for i := range buf {
		buf[i] = 0
	}
	// Reset length to 0
	buf = buf[:0]
	runeBufferPool.Put(buf)
}
