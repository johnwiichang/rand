# rand

**Rand** package generates fast, cryptographic random numbers using the [RDRAND](https://en.wikipedia.org/wiki/RdRand) instruction. Some argue that RDRAND is cryptographically compromised by National Security Agency using an unknown backdoor. On the other hand, using the instruction is ~3x faster than *crypto/rand* (on a Mac using go test benchmark).

The package accesses assembly instructions via C-go C wrapper.

## Example
    package main

    import (
        log "github.com/surendarchandra/rand"
    )

    func main() {
        if rand.IsRand64Supported() {
            var rnd uint64
            RdRand64(&rnd)

            /* Byte arrays are also supported */
            b = make([]byte, 1000)
            rand.Read(b)
        }
    }

## Performance Results
### Mac OSX
* crypto/rand, 1000 bytes: 63779 ns/op
* crypto/rand, 8 bytes: 592   ns/op
* math/rand, 8 bytes: 6.10 ns/op
* this package, 1000 bytes:19616    ns/op
* this package, 8 bytes	: 159    ns/op
