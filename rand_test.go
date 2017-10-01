package rand

import (
	crand "crypto/rand"
	mrand "math/rand"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type RandSuite struct{}

var _ = Suite(&RandSuite{})

// Random numbers are generated as 64bit values. Check corner cases
func (r *RandSuite) TestRdRand(c *C) {
	var rnd uint64

	if IsRand64Supported() {
		RdRand64(&rnd)

		b := make([]byte, 5)
		Read(b)

		b = make([]byte, 8)
		Read(b)

		b = make([]byte, 1000)
		Read(b)
	}

	if IsSeed64Supported() {
		// RDSEED requires Broadwell or newer
		RdSeed64(&rnd)
	}
}

// math random routines are not cryptographic strength but fast
func (r *RandSuite) BenchmarkMathRand8(c *C) {
	rnd := mrand.New(mrand.NewSource(0))

	for i := 0; i < c.N; i++ {
		rnd.Uint64()
	}

}

func (r *RandSuite) BenchmarkRdRand8(c *C) {
	if !IsRand64Supported() {
		return
	}

	b := make([]byte, 8)
	for i := 0; i < c.N; i++ {
		Read(b)
	}
}

func (r *RandSuite) BenchmarkRdRand1000(c *C) {
	if !IsRand64Supported() {
		return
	}

	b := make([]byte, 1000)
	for i := 0; i < c.N; i++ {
		Read(b)
	}
}

func (r *RandSuite) BenchmarkCryptoRand8(c *C) {
	b := make([]byte, 8)
	for i := 0; i < c.N; i++ {
		crand.Read(b)
	}

}

func (r *RandSuite) BenchmarkCryptoRand1000(c *C) {
	b := make([]byte, 1000)
	for i := 0; i < c.N; i++ {
		crand.Read(b)
	}

}
