package randInit

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
)

func init() {
	n, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(n.Int64())
}
