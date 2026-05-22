package pocsag

import (
	"errors"
	"math/bits"
)

const (
	bchGenerator  uint32 = 0x769 // 11101101001
	bchN          uint32 = 31
	bchK          uint32 = 21
	bchParityBits uint32 = 10
)

const syncCodeword uint32 = 0x7CD215D8
const idleCodeword uint32 = 0x7A89C197

var errUncorrectable = errors.New("bch: uncorrectable error (>2 bits)")

func bchEncode(data uint32) uint32 {
	shifted := (data & 0x1FFFFF) << bchParityBits
	remainder := shifted
	for i := bchN - 1; i >= bchParityBits; i-- {
		if remainder&(1<<i) != 0 {
			remainder ^= bchGenerator << (i - bchParityBits)
		}
	}
	return shifted | (remainder & 0x3FF)
}

func bchDecode(codeword uint32) (uint32, error) {
	syndrome := calcSyndrome(codeword)
	if syndrome == 0 {
		return (codeword >> bchParityBits) & 0x1FFFFF, nil
	}
	// Check single-bit errors
	for i := uint32(0); i < bchN; i++ {
		errPattern := uint32(1) << i
		if calcSyndrome(codeword^errPattern) == 0 {
			corrected := codeword ^ errPattern
			return (corrected >> bchParityBits) & 0x1FFFFF, nil
		}
	}
	// Check two-bit errors
	for i := uint32(0); i < bchN-1; i++ {
		for j := i + 1; j < bchN; j++ {
			errPattern := (uint32(1) << i) | (uint32(1) << j)
			if calcSyndrome(codeword^errPattern) == 0 {
				corrected := codeword ^ errPattern
				return (corrected >> bchParityBits) & 0x1FFFFF, nil
			}
		}
	}
	return 0, errUncorrectable
}

func calcSyndrome(cw uint32) uint32 {
	remainder := cw & 0x7FFFFFFF
	for i := bchN - 1; i >= bchParityBits; i-- {
		if remainder&(1<<i) != 0 {
			remainder ^= bchGenerator << (i - bchParityBits)
		}
	}
	return remainder & 0x3FF
}

type frameSynchronizer struct {
	state    int
	buffer   [32]byte
	bufPos   int
	bitCount int
}

func newFrameSynchronizer() *frameSynchronizer {
	return &frameSynchronizer{}
}

const maxSyncErrors = 2

func hammingDistance(a, b uint32) int {
	return bits.OnesCount32(a ^ b)
}

func (fs *frameSynchronizer) processBits(bits []byte) []uint32 {
	var results []uint32
	for _, b := range bits {
		fs.buffer[fs.bufPos] = b
		fs.bufPos = (fs.bufPos + 1) % 32
		fs.bitCount++
		if fs.bitCount%32 == 0 {
			var cw uint32
			for i := 0; i < 32; i++ {
				bitIdx := (fs.bufPos - 32 + i + 32) % 32
				if fs.buffer[bitIdx] == 1 {
					cw |= 1 << (31 - i)
				}
			}
			if fs.state == 0 {
				if hammingDistance(cw, syncCodeword) <= maxSyncErrors {
					fs.state = 1
				}
			} else {
				results = append(results, cw)
			}
		}
	}
	return results
}

func generate1010Bits(n int) []byte {
	bits := make([]byte, n)
	for i := range bits {
		bits[i] = byte(i & 1)
	}
	return bits
}

func uint32ToBits(v uint32) []byte {
	bits := make([]byte, 32)
	for i := 0; i < 32; i++ {
		if v&(1<<(31-i)) != 0 {
			bits[i] = 1
		}
	}
	return bits
}
