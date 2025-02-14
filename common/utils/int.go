package utils

import "math"

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func ToInt64[N Integer](n N) int64 {
	if n > 0 && uint64(n) > math.MaxInt64 {
		panic("Integer overflow")
	}
	return int64(n)
}

func ToInt32[N Integer](n N) int32 {
	if (n > 0 && uint64(n) > math.MaxInt32) || (n < 0 && int64(n) < math.MinInt32) {
		panic("Integer overflow")
	}
	return int32(n)
}

func ToUint64[N Integer](n N) uint64 {
	if n < 0 {
		panic("Integer overflow: must not be negative")
	}
	return uint64(n)
}

// func ToUint32[N Integer](n N) uint32 {
// 	if n < 0 || uint64(n) > math.MaxUint32 {
// 		panic("Integer overflow")
// 	}
// 	return uint32(n)
// }

// func ToUint16[N Integer](n N) uint16 {
// 	if n < 0 || uint64(n) > math.MaxUint16 {
// 		panic("Integer overflow")
// 	}
// 	return uint16(n)
// }
