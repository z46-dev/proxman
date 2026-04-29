package tests

import (
	"math/rand/v2"
	"testing"
)

func randomBytes(n int) (b []byte) {
	b = make([]byte, n)
	for i := range b {
		b[i] = byte(rand.IntN(256))
	}

	return
}

func mustConfig(t *testing.T) {
	var err error

	if err = initConfig("testconfig.toml"); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
}
