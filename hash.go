package bitacoin

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// GenerateMask creates a mask based on the number of zeros required in the hash
func GenerateMask(zeros int) []byte {
	full, half := zeros/2, zeros%2
	var mask []byte
	for i := 0; i < full; i++ {
		mask = append(mask, 0)
	}

	if half > 0 {
		mask = append(mask, 0xf)
	}

	return mask
}

// GoodEnough checks if the hash is good for the current mask
func GoodEnough(mask []byte, hash []byte) bool {
	for i := range mask {
		if hash[i] > mask[i] {
			return false
		}
	}
	return true
}

//CompareHash check if nonce is acceptable with difficulty of hash
func CompareHash(difficulty int, hash []byte) bool {
	var hashInt big.Int
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	hashInt.SetBytes(hash[:])

	if hashInt.Cmp(target) == -1 {
		return true
	}
	return false
}

// EasyHash craete hash, the easy way, just a simple sha256 hash
func EasyHash(data ...interface{}) []byte {
	hasher := sha256.New()

	fmt.Fprint(hasher, data...)

	return hasher.Sum(nil)
}

// DifficultHash creates the hash with difficulty mask and conditions,
// return the hash and the nonce used to create the hash
func DifficultHash(difficulty int, data ...interface{}) ([]byte, int32) {
	ln := len(data)
	data = append(data, nil)
	var i int32
	for {
		data[ln] = i
		hash := EasyHash(data...)
		if CompareHash(difficulty, hash) {
			return hash, i
		}
		i++
	}
}
