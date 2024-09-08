package hasher

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func toBase62(num uint64) []byte {
	var result [6]byte
	for i := uint(0); i < 6; i++ {
		result[5-i] = base62Chars[num%62]
		num /= 62
	}

	return result[:]
}

func GetUrlHash(url string) []byte {
	randSalt := rand.Int63n(1000000)
	input := fmt.Sprintf("%s%d", url, randSalt)

	hash := sha256.Sum256([]byte(input))
	hashPrefix := binary.BigEndian.Uint64(hash[:8])

	return toBase62(hashPrefix)
}
