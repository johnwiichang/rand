package rand

import (
	"errors"
	"unsafe"

	"github.com/klauspost/cpuid"
)

//+build !386

// #cgo amd64
// Based on Intel® Digital Random Number Generator (DRNG) Software Implementation Guide
//   By John M. (Intel)
// Supports 64 bit random numbers; will loop indefinitely
// until enough entropy is accumulated

// unsigned char sizeMismatch(unsigned long long *val) {
//   return sizeof(*val) == 8 ? 1 : 0;
// }
//
// void rdrand64(unsigned long long *rand) {
//   unsigned char ok;
//
//   while (1) {
//     // Loop until enough entropy
//     asm volatile(".byte 0x48,0x0f,0xc7,0xf0; setc %1"
//                : "=a" (*rand), "=qm" (ok)
//                :
//                : "rdx"
//               );
//     if (ok == 1) {
//       return;
//     }
//   }
// }
//
// void rdseed64(unsigned long long *seed) {
//   unsigned char ok;
//
//   while (1) {
//     asm volatile(".byte 0x48,0x0f,0xc7,0xf8; setc %1"
//              : "=a" (*seed), "=qm" (ok)
//              :
//              : "rdx"
//              );
//     if (ok == 1) {
//       return;
//     }
//   }
// }
import "C"

const (
	maxBufSize = 1000000
)

var rdRandSupported, rdSeedSupported bool

func init() {
	var val C.ulonglong
	if C.sizeMismatch(&val) == 0 {
		// C unsigned long long is not the same as GO uint64
		return
	}
	rdRandSupported = cpuid.CPU.Rdrand()
	rdSeedSupported = cpuid.CPU.Rdseed()
}

// IsRand64Supported reports whether RDRAND is supported. Caller is
// responsible for first checking before issuing RdRand64
func IsRand64Supported() bool {
	return rdRandSupported
}

// IsSeed64Supported reports whether RDSEEDD is supported. Caller is
// responsible for first checking before issuing RdSeed64
func IsSeed64Supported() bool {
	return rdSeedSupported
}

// RdRand64 generates a 64bit random number. Assumes that RDRAND is supported.
func RdRand64(rand *uint64) {
	C.rdrand64((*C.ulonglong)(rand))
}

// RdSeed64 seeds a 64bit value. Assumes that RDSEED is supported
func RdSeed64(seed *uint64) {
	C.rdseed64((*C.ulonglong)(seed))
}

// Read fills the byte array with random values. Returns the size of buffer when no error.
func Read(b []byte) (n int, err error) {
	size := len(b)
	size64 := size / 8
	if size64 > maxBufSize {
		return 0, errors.New("Too big")
	}

	// Convert the []byte into a uint64 array, unsafely.
	p64 := (*[maxBufSize]uint64)(unsafe.Pointer(&b[0]))
	for i := 0; i < size64; i++ {
		RdRand64(&p64[i])
	}

	leftover := size - size64*8
	if leftover == 0 {
		return size, nil
	}

	// Fill the remaining bytes
	var rnd uint64
	RdRand64(&rnd)

	r8 := (*[8]uint8)(unsafe.Pointer(&rnd))
	copy(b[size64*8:], r8[:leftover])

	return size, nil
}
