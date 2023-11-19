package pow

import (
	"crypto/sha1"
	"fmt"
)

const zeroByte = 48

// HashcashData challenge structure
type HashcashData struct {
	Version    int    // Hashcash format version
	ZerosCount int    // Number of zero bits in the hashed code.
	Date       int64  // The time that the message was sent
	Resource   string // Resource data string being transmitted (Client IP will be used for this project)
	Rand       string // String of random characters, encoded in base-64 format.
	Counter    int    // Binary counter, encoded in base-64 format
}

func (h HashcashData) Stringify() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter)
}

func sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// IsHashCorrect checks if hash has zerosCount zeros at the beginning
func IsHashCorrect(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}
	for _, ch := range hash[:zerosCount] {
		if ch != zeroByte {
			return false
		}
	}
	return true
}

// ComputeHashcash generates the correct hashcash by continuously attempting
// brute force calculations until the resulting hash meets the
// condition of IsHashCorrect
func (h HashcashData) ComputeHashcash(maxIterations int) (HashcashData, error) {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		header := h.Stringify()
		hash := sha1Hash(header)
		if IsHashCorrect(hash, h.ZerosCount) {
			return h, nil
		}
		h.Counter++
	}
	return h, fmt.Errorf("max iterations exceeded")
}
