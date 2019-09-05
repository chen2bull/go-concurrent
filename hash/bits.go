package hash

import "math/bits"

var tableBitReverse = [...]int{
	0x00, 0x80, 0x40, 0xC0, 0x20, 0xA0, 0x60, 0xE0, 0x10, 0x90, 0x50, 0xD0, 0x30, 0xB0, 0x70, 0xF0,
	0x08, 0x88, 0x48, 0xC8, 0x28, 0xA8, 0x68, 0xE8, 0x18, 0x98, 0x58, 0xD8, 0x38, 0xB8, 0x78, 0xF8,
	0x04, 0x84, 0x44, 0xC4, 0x24, 0xA4, 0x64, 0xE4, 0x14, 0x94, 0x54, 0xD4, 0x34, 0xB4, 0x74, 0xF4,
	0x0C, 0x8C, 0x4C, 0xCC, 0x2C, 0xAC, 0x6C, 0xEC, 0x1C, 0x9C, 0x5C, 0xDC, 0x3C, 0xBC, 0x7C, 0xFC,
	0x02, 0x82, 0x42, 0xC2, 0x22, 0xA2, 0x62, 0xE2, 0x12, 0x92, 0x52, 0xD2, 0x32, 0xB2, 0x72, 0xF2,
	0x0A, 0x8A, 0x4A, 0xCA, 0x2A, 0xAA, 0x6A, 0xEA, 0x1A, 0x9A, 0x5A, 0xDA, 0x3A, 0xBA, 0x7A, 0xFA,
	0x06, 0x86, 0x46, 0xC6, 0x26, 0xA6, 0x66, 0xE6, 0x16, 0x96, 0x56, 0xD6, 0x36, 0xB6, 0x76, 0xF6,
	0x0E, 0x8E, 0x4E, 0xCE, 0x2E, 0xAE, 0x6E, 0xEE, 0x1E, 0x9E, 0x5E, 0xDE, 0x3E, 0xBE, 0x7E, 0xFE,
	0x01, 0x81, 0x41, 0xC1, 0x21, 0xA1, 0x61, 0xE1, 0x11, 0x91, 0x51, 0xD1, 0x31, 0xB1, 0x71, 0xF1,
	0x09, 0x89, 0x49, 0xC9, 0x29, 0xA9, 0x69, 0xE9, 0x19, 0x99, 0x59, 0xD9, 0x39, 0xB9, 0x79, 0xF9,
	0x05, 0x85, 0x45, 0xC5, 0x25, 0xA5, 0x65, 0xE5, 0x15, 0x95, 0x55, 0xD5, 0x35, 0xB5, 0x75, 0xF5,
	0x0D, 0x8D, 0x4D, 0xCD, 0x2D, 0xAD, 0x6D, 0xED, 0x1D, 0x9D, 0x5D, 0xDD, 0x3D, 0xBD, 0x7D, 0xFD,
	0x03, 0x83, 0x43, 0xC3, 0x23, 0xA3, 0x63, 0xE3, 0x13, 0x93, 0x53, 0xD3, 0x33, 0xB3, 0x73, 0xF3,
	0x0B, 0x8B, 0x4B, 0xCB, 0x2B, 0xAB, 0x6B, 0xEB, 0x1B, 0x9B, 0x5B, 0xDB, 0x3B, 0xBB, 0x7B, 0xFB,
	0x07, 0x87, 0x47, 0xC7, 0x27, 0xA7, 0x67, 0xE7, 0x17, 0x97, 0x57, 0xD7, 0x37, 0xB7, 0x77, 0xF7,
	0x0F, 0x8F, 0x4F, 0xCF, 0x2F, 0xAF, 0x6F, 0xEF, 0x1F, 0x9F, 0x5F, 0xDF, 0x3F, 0xBF, 0x7F, 0xFF,
}

const (
	// 哨兵元素的哈希值 h & mask以后，hiMask32所在的位必定为0，再reverse32以后，最低位必定为0
	// 普通元素的哈希值 (h & mask)|hiMask32以后，hiMask32所在的位必定为1，再reverse32以后，最低位必定为1
	mask32     int32 = 0x3FFFFFFF
	wordSize32 int32 = 31 // 减去符号位不用
	hiMask32   int32 = 0x40000000
	loMask32   int32 = 0x00000001
	// maxCap32值为1 << 29，bucket数组的最大容量，可以保证最大下标的hiMask32位必然为0
	// reverse32以后最低位必定为0（即保证了作为哨兵节点的key最低位必为0）
	maxBucketSize32 int32 = 0x20000000
	minBucketSize32 int32 = 32 // Must be a32 power of 2

	// 哨兵元素的哈希值 h & mask以后，hiMask64所在的位必定为0，再reverse64以后，最低位必定为0
	// 普通元素的哈希值 (h & mask)|hiMask64以后，hiMask64所在的位必定为1，再reverse64以后，最低位必定为1
	mask64          int64 = 0x007FFFFFFFFFFFFF
	wordSize64      int64 = 56
	loMask64        int64 = 0x0000000000000001
	hiMask64        int64 = 0x0080000000000000
	maxBucketSize64 int64 = 0x0040000000000000
	minBucketSize64 int64 = 32
	// maxCap64值为1 << 54，bucket数组的最大容量，可以保证最大下标的HiMask必然为0
	// reverse64以后最低位必定为0（即保证了作为哨兵节点的key最低位必为0）

	DefaultLoadFactor = 4 // float64(mapSizeNow)/float64(bucketSizeNow)  > DefaultLoadFactor的时候，扩展桶数组
)

var mask int
var hiMask int

func makeRegularKey(hashKey int64) int64 {
	code := hashKey & mask64
	return lookupReverse64(code | hiMask64)
}

func isRegularKey(hashKey int64) bool {
	return (hashKey & loMask64) == 1
}

func makeSentinelKey(hashKey int64) int64 {
	return lookupReverse64(hashKey & mask64)
}

// reverse前必须已经 & mask32
func reverse32(key int32) int32 {
	loMask := loMask32
	hiMask := hiMask32
	result := int32(0)
	for i := int32(0); i < wordSize32; i++ {
		if (key & loMask) != 0 { // bit set
			result |= hiMask
		}
		loMask <<= 1
		hiMask = hiMask >> 1 // fill with 0 from left
	}
	return result
}

// reverse前必须已经 & mask64
func reverse64(key int64) int64 {
	loMask := loMask64
	hiMask := hiMask64
	result := int64(0)
	for i := int64(0); i < wordSize64; i++ {
		if (key & loMask) != 0 { // bit set
			result |= hiMask
		}
		loMask <<= 1
		hiMask = hiMask >> 1 // fill with 0 from left
	}
	return result
}

var reverse func(key int) int

func lookupReverse32(key int32) int32 {
	result := int32(0)
	result = result | int32(tableBitReverse[key&0xff]<<23)      // 23==16 + 7
	result = result | int32(tableBitReverse[(key>>8)&0xff]<<15) // 15 == 8 + 7
	result = result | int32(tableBitReverse[(key>>16)&0xff]<<7)
	result = result | int32(tableBitReverse[((key>>24)<<1)&0xff]) // 右移完再左移一位，最低位肯定为0，翻转后就只有最低位
	return result
}

func lookupReverse64(key int64) int64 {
	result := int64(0)
	result = result | int64(tableBitReverse[key&0xff]<<48)
	result = result | int64(tableBitReverse[(key>>8)&0xff]<<40)
	result = result | int64(tableBitReverse[(key>>16)&0xff]<<32)
	result = result | int64(tableBitReverse[(key>>24)&0xff]<<24)
	result = result | int64(tableBitReverse[(key>>32)&0xff]<<16)
	result = result | int64(tableBitReverse[(key>>40)&0xff]<<8)
	result = result | int64(tableBitReverse[(key>>48)&0xff])
	return result
}

var lookupReverse func(key int) int

func tableSizeFor32(c int32) int32 {
	if c <= minBucketSize32 {
		return minBucketSize32
	}
	if c >= maxBucketSize32 {
		return maxBucketSize32
	}
	n := uint32(c - 1)
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n+1 >= uint32(maxBucketSize32) {
		return maxBucketSize32
	}
	return int32(n + 1)
}

func tableSizeFor64(c int64) int64 {
	if c <= minBucketSize64 {
		return minBucketSize64
	}
	if c >= maxBucketSize64 {
		return maxBucketSize64
	}
	n := uint64(c - 1)
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	if n+1 >= uint64(maxBucketSize64) {
		return maxBucketSize64
	}
	return int64(n + 1)
}

var tableSizeFor func(key int) int

func isPowerOfTwo64(x int64) bool {
	return (x != 0) && ((x & (^x + 1)) == x)
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func init() {
	if bits.UintSize == 32 {
		reverse = func(key int) int {
			return int(reverse32(int32(key)))
		}
		lookupReverse = func(key int) int {
			return int(lookupReverse32(int32(key)))
		}
		tableSizeFor = func(key int) int {
			return int(tableSizeFor32(int32(key)))
		}
		mask = int(mask32)
		hiMask = int(hiMask32)
	} else {
		reverse = func(key int) int {
			return int(reverse64(int64(key)))
		}
		lookupReverse = func(key int) int {
			return int(lookupReverse64(int64(key)))
		}
		tableSizeFor = func(key int) int {
			return int(tableSizeFor64(int64(key)))
		}
		mask = int(mask64)
		hiMask = int(hiMask64)
	}
}
