crc [![GoDoc](https://godoc.org/github.com/snksoft/src?status.png)](https://godoc.org/github.com/snksoft/crc)
========
This package implements generic CRC calculations up to 64 bits wide.
It aims to be fairly fast and fairly complete, allowing users to match pretty much
any CRC algorithm used in the wild by choosing appropriate Parameters. This obviously includes all popular CRC algorithms, such as CRC64-ISO, CRC64-ECMA, CRC32, CRC32C, CRC16, CCITT, XMODEM and many others. See http://reveng.sourceforge.net/crc-catalogue/ for a good list of CRC algorithms and their parameters.

This package has been largely inspired by Ross Williams' 1993 paper "A Painless Guide to CRC Error Detection Algorithms".


## Installation

To install, simply execute:

```
go get github.com/snksoft/crc
```

## Usage

Using `crc` is easy. Here is an example of calculating a CCITT crc.
```go
package main

import (
	"fmt"
	"github.com/snksoft/crc"
)

func main() {
	data := "123456789"
	ccittCrc := crc.CalculateCRC(crc.CCITT, []byte(data))
	fmt.Printf("CRC is 0x%04X\n", ccittCrc) // prints "CRC is 0x29B1"
}
```

For larger data, table driven implementation is faster. Note that `crc.Hash` implements `hash.Hash` interface, so you can use it instead if you want.  
Here is how to use it:
```go
package main

import (
	"fmt"
	"github.com/snksoft/crc"
)

func main() {
	data := "123456789"
	hash := crc.NewHash(crc.XMODEM)
	xmodemCrc := hash.CalculateCRC([]byte(data))
	fmt.Printf("CRC is 0x%04X\n", xmodemCrc) // prints "CRC is 0x31C3"

	// You can also reuse hash instance for another crc calculation
	// And if data is too big, you may feed it in chunks
	hash.Reset() // Discard crc data accumulated so far
	hash.Update([]byte("123456789")) // feed first chunk
	hash.Update([]byte("01234567890")) // feed next chunk
	xmodemCrc2 := hash.CRC() // gets CRC of whole data ("12345678901234567890")
	fmt.Printf("CRC is 0x%04X\n", xmodemCrc2) // prints "CRC is 0x2C89"
}
```
## Notes
Beware that Hash instance is not thread safe. If you want to do parallel CRC calculations, then either use `NewHash()` to create multiple Hash instances or simply make a copy of Hash whehever you need it. Latter option avoids recalculating CRC table, but keep in mind that `NewHash()` returns a pointer, so simple assignement will point to the same instance.
Use either
 ```go
hash2 := &crc.Hash{}
*hash2 = *hash
```
or simply
 ```go
var hash2 = *hash
 ```
