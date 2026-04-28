package tests

import "math/rand/v2"

func randomBytes(n int) (b []byte) {
	b = make([]byte, n)
	for i := range b {
		b[i] = byte(rand.IntN(256))
	}

	return
}
