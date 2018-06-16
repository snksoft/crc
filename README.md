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

or:

```
go get gopkg.in/snksoft/crc.v1
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
	xmodemCrc2 := hash.CRC() // gets CRC of whole data "12345678901234567890"
	fmt.Printf("CRC is 0x%04X\n", xmodemCrc2) // prints "CRC is 0x2C89"
}
```

## New in version 1.1

In this version I have separated actual CRC caclulations and Hash interface implementation. New `Table` type incorporates table based implementation which can be used without creating a `Hash` instance. The main difference is that `Table` instances are essentially immutable once initialized. This greatly simplifies concurrent use as `Table` instances can be safely used in concurrent applications without tricky copying or synchronization. The downside is, however, that feeding data in multiple chunks becomes a bit more verbose (as you essentially maintain intermediate crc in your code and keep feeding it back to subsequent calls). So, you might prefer one or the other depending on situation at hand and personal preferences. You even can ask a `Hash` instance for a `Table` instance it uses internally and then use both in parallel without recalculating the crc table.

Anyway, here is how to use a `Table` directly.

```go
package main

import (
	"fmt"
	"github.com/snksoft/crc"
)

func main() {
	data := []byte("123456789")

	// create a Table
	crcTable := crc.NewTable(crc.XMODEM)

	// Simple calculation all in one go
	xmodemCrc := crcTable.CalculateCRC(data)
	fmt.Printf("CRC is 0x%04X\n", xmodemCrc) // prints "CRC is 0x31C3"

	// You can also reuse same Table for another crc calculation
	// or even calculate multiple crc in parallel using same Table
	crc1 := crcTable.InitCrc()
	crc1 = crcTable.UpdateCrc(crc1, []byte("1234567890")) // feed first chunk to first crc
	crc2 := crcTable.InitCrc()
	crc2 = crcTable.UpdateCrc(crc2, data)                  // feed first chunk to second crc
	crc1 = crcTable.UpdateCrc(crc1, []byte("1234567890")) // feed second chunk to first crc

	// Now finish calcuation for both
	crc1 = crcTable.CRC(crc1)
	crc2 = crcTable.CRC(crc2)

	fmt.Printf("CRC is 0x%04X\n", crc1) // prints "CRC is 0x2C89"
	fmt.Printf("CRC is 0x%04X\n", crc2) // prints "CRC is 0x31C3"
}
```


## Notes
Beware that `Hash` instance is not thread safe. If you want to do parallel CRC calculations (and actually need it to be `Hash`, not `Table`), then either use `NewHash()` to create multiple Hash instances or simply make a copy of Hash whehever you need it. Latter option avoids recalculating CRC table, but keep in mind that `NewHash()` returns a pointer, so simple assignement will point to the same instance.
Use either
 ```go
hash2 := &crc.Hash{}
*hash2 = *hash
```
or simply
 ```go
var hash2 = *hash
 ```
