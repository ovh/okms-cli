package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToInt64(t *testing.T) {
	assert.Equal(t, int64(12), ToInt64(int8(12)))
	assert.Equal(t, int64(-12), ToInt64(int8(-12)))
	assert.Equal(t, int64(12), ToInt64(uint8(12)))
	assert.Equal(t, int64(12), ToInt64(int64(12)))
	assert.Equal(t, int64(12), ToInt64(uint64(12)))
	assert.Equal(t, int64(0), ToInt64(int64(0)))
	assert.Equal(t, int64(0), ToInt64(uint64(0)))

	assert.Panics(t, func() {
		ToInt64(uint64(math.MaxInt64 + 1))
	})
}

func TestToInt32(t *testing.T) {
	assert.Equal(t, int32(12), ToInt32(int8(12)))
	assert.Equal(t, int32(-12), ToInt32(int8(-12)))
	assert.Equal(t, int32(12), ToInt32(uint8(12)))
	assert.Equal(t, int32(12), ToInt32(int64(12)))
	assert.Equal(t, int32(12), ToInt32(uint64(12)))
	assert.Equal(t, int32(0), ToInt32(int64(0)))
	assert.Equal(t, int32(0), ToInt32(uint64(0)))

	assert.Panics(t, func() {
		ToInt32(uint32(math.MaxInt32 + 1))
	})

	assert.Panics(t, func() {
		ToInt32(int64(math.MinInt32 - 1))
	})
}

// func TestToUint32(t *testing.T) {
// 	assert.Equal(t, uint32(12), ToUint32(int8(12)))
// 	assert.Equal(t, uint32(12), ToUint32(uint8(12)))
// 	assert.Equal(t, uint32(12), ToUint32(int64(12)))
// 	assert.Equal(t, uint32(12), ToUint32(uint64(12)))
// 	assert.Equal(t, uint32(0), ToUint32(int64(0)))
// 	assert.Equal(t, uint32(0), ToUint32(uint64(0)))

// 	assert.Panics(t, func() {
// 		ToUint32(uint64(math.MaxUint32 + 1))
// 	})

// 	assert.Panics(t, func() {
// 		ToUint32(int8(-12))
// 	})
// }

// func TestToUint16(t *testing.T) {
// 	assert.Equal(t, uint16(12), ToUint16(int8(12)))
// 	assert.Equal(t, uint16(12), ToUint16(uint8(12)))
// 	assert.Equal(t, uint16(12), ToUint16(int64(12)))
// 	assert.Equal(t, uint16(12), ToUint16(uint64(12)))
// 	assert.Equal(t, uint16(0), ToUint16(int64(0)))
// 	assert.Equal(t, uint16(0), ToUint16(uint64(0)))

// 	assert.Panics(t, func() {
// 		ToUint16(uint64(math.MaxUint16 + 1))
// 	})

// 	assert.Panics(t, func() {
// 		ToUint16(int8(-12))
// 	})
// }
