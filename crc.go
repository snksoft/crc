// Copyright 2016, S&K Software Development Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crc implements generic CRC calculations up to 64 bits wide.
// It aims to be fairly complete, allowing users to match pretty much
// any CRC algorithm used in the wild by choosing appropriate Parameters.
// And it's also fairly fast for everyday use.
//
// This package has been largely inspired by Ross Williams' 1993 paper "A Painless Guide to CRC Error Detection Algorithms".
// A good list of parameter sets for various CRC algorithms can be found at http://reveng.sourceforge.net/crc-catalogue/.
package crc

// Parameters represents set of parameters defining a particular CRC algorithm.
type Parameters struct {
	Width      uint   // Width of the CRC expressed in bits
	Polynomial uint64 // Polynomial used in this CRC calculation
	ReflectIn  bool   // ReflectIn indicates whether input bytes should be reflected
	ReflectOut bool   // ReflectOut indicates whether input bytes should be reflected
	Init       uint64 // Init is initial value for CRC calculation
	FinalXor   uint64 // Xor is a value for final xor to be applied before returning result
}

var (
	// X-25 CRC parameters, also known as CRC-16/IBM-SDLC, CRC-16/ISO-HDLC, CRC-B
	X25 = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0xFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFF}
	// CCITT CRC parameters
	CCITT = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0xFFFF, ReflectIn: false, ReflectOut: false, FinalXor: 0x0}
	// CRC16 CRC parameters, also known as ARC
	CRC16 = &Parameters{Width: 16, Polynomial: 0x8005, Init: 0x0000, ReflectIn: true, ReflectOut: true, FinalXor: 0x0}
	// XMODEM is a set of CRC parameters commonly referred as "XMODEM"
	XMODEM = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0}
	// XMODEM2 is another set of CRC parameters commonly referred as "XMODEM"
	XMODEM2 = &Parameters{Width: 16, Polynomial: 0x8408, Init: 0x0000, ReflectIn: true, ReflectOut: true, FinalXor: 0x0}

	// CRC32 is by far the the most commonly used CRC-32 polynom and set of parameters
	CRC32 = &Parameters{Width: 32, Polynomial: 0x04C11DB7, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFF}
	// IEEE is an alias to CRC32
	IEEE = CRC32
	// Castagnoli polynomial. used in iSCSI. And also provided by hash/crc32 package.
	Castagnoli = &Parameters{Width: 32, Polynomial: 0x1EDC6F41, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFF}
	// CRC32C is an alias to Castagnoli
	CRC32C = Castagnoli
	// Koopman polynomial
	Koopman = &Parameters{Width: 32, Polynomial: 0x741B8CD7, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFF}

	// CRC64ISO is set of parameters commonly known as CRC64-ISO
	CRC64ISO = &Parameters{Width: 64, Polynomial: 0x000000000000001B, Init: 0xFFFFFFFFFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFFFFFFFFFF}
	// CRC64ECMA is set of parameters commonly known as CRC64-ECMA
	CRC64ECMA = &Parameters{Width: 64, Polynomial: 0x42F0E1EBA9EA3693, Init: 0xFFFFFFFFFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFFFFFFFFFF}
)

// reflect reverses order of last count bits
func reflect(in uint64, count uint) uint64 {
	ret := in
	for idx := uint(0); idx < count; idx++ {
		srcbit := uint64(1) << idx
		dstbit := uint64(1) << (count - idx - 1)
		if (in & srcbit) != 0 {
			ret |= dstbit
		} else {
			ret = ret & (^dstbit)
		}
	}
	return ret
}

// CalculateCRC implements simple straight forward bit by bit calculation.
// It is relatively slow for large amounts of data, but does not require
// any preparation steps. As a result, it might be faster in some cases
// then building a table required for faster calculation.
//
// Note: this implementation follows section 8 ("A Straightforward CRC Implementation")
// of Ross N. Williams paper as even though final/sample implementation of this algorithm
// provided near the end of that paper (and followed by most other implementations)
// is a bit faster, it does not work for polynomials shorter then 8 bits. And if you need
// speed, you shoud probably be using table based implementation anyway.
func CalculateCRC(crcParams *Parameters, data []byte) uint64 {

	curValue := crcParams.Init
	topBit := uint64(1) << (crcParams.Width - 1)
	mask := (topBit << 1) - 1

	for i := 0; i < len(data); i++ {
		var curByte = uint64(data[i]) & 0x00FF
		if crcParams.ReflectIn {
			curByte = reflect(curByte, 8)
		}
		for j := uint64(0x0080); j != 0; j >>= 1 {
			bit := curValue & topBit
			curValue <<= 1
			if (curByte & j) != 0 {
				bit = bit ^ topBit
			}
			if bit != 0 {
				curValue = curValue ^ crcParams.Polynomial
			}
		}
	}
	if crcParams.ReflectOut {
		curValue = reflect(curValue, crcParams.Width)
	}

	curValue = curValue ^ crcParams.FinalXor

	return curValue & mask
}

// Table represents the partial evaluation of a checksum using table-driven
// implementation. It is essentially immutable once initialized and thread safe as a result.
type Table struct {
	crcParams Parameters
	crctable  []uint64
	mask      uint64
	initValue uint64
}

// NewTable creates and initializes a new Table for the CRC algorithm specified by the crcParams.
func NewTable(crcParams *Parameters) *Table {
	ret := &Table{crcParams: *crcParams}
	ret.mask = (uint64(1) << crcParams.Width) - 1
	ret.crctable = make([]uint64, 256, 256)
	ret.initValue = crcParams.Init
	if crcParams.ReflectIn {
		ret.initValue = reflect(crcParams.Init, crcParams.Width)
	}

	tmp := make([]byte, 1, 1)
	tableParams := *crcParams
	tableParams.Init = 0
	tableParams.ReflectOut = tableParams.ReflectIn
	tableParams.FinalXor = 0
	for i := 0; i < 256; i++ {
		tmp[0] = byte(i)
		ret.crctable[i] = CalculateCRC(&tableParams, tmp)
	}
	return ret
}

// InitCrc returns a stating value for a new CRC calculation
func (t *Table) InitCrc() uint64 {
	return t.initValue
}

// UpdateCrc process supplied bytes and updates current (partial) CRC accordingly.
// It can be called repetitively to process larger data in chunks.
func (t *Table) UpdateCrc(curValue uint64, p []byte) uint64 {
	if t.crcParams.ReflectIn {
		for _, v := range p {
			curValue = t.crctable[(byte(curValue)^v)&0xFF] ^ (curValue >> 8)
		}
	} else if t.crcParams.Width < 8 {
		for _, v := range p {
			curValue = t.crctable[((((byte)(curValue<<(8-t.crcParams.Width)))^v)&0xFF)] ^ (curValue << 8)
		}
	} else {
		for _, v := range p {
			curValue = t.crctable[((byte(curValue>>(t.crcParams.Width-8))^v)&0xFF)] ^ (curValue << 8)
		}
	}
	return curValue
}

// CRC returns CRC value for the data processed so far.
func (t *Table) CRC(curValue uint64) uint64 {
	ret := curValue

	if t.crcParams.ReflectOut != t.crcParams.ReflectIn {
		ret = reflect(ret, t.crcParams.Width)
	}
	return (ret ^ t.crcParams.FinalXor) & t.mask
}

// CRC8 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (t *Table) CRC8(curValue uint64) uint8 {
	return uint8(t.CRC(curValue))
}

// CRC16 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (t *Table) CRC16(curValue uint64) uint16 {
	return uint16(t.CRC(curValue))
}

// CRC32 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (t *Table) CRC32(curValue uint64) uint32 {
	return uint32(t.CRC(curValue))
}

// CalculateCRC is a convenience function allowing to calculate CRC in one call.
func (t *Table) CalculateCRC(data []byte) uint64 {
	crc := t.InitCrc()
	crc = t.UpdateCrc(crc, data)
	return t.CRC(crc)
}

// Hash represents the partial evaluation of a checksum using table-driven
// implementation. It also implements hash.Hash interface.
type Hash struct {
	table    *Table
	curValue uint64
	size     uint
}

// Size returns the number of bytes Sum will return.
// See hash.Hash interface.
func (h *Hash) Size() int { return int(h.size) }

// BlockSize returns the hash's underlying block size.
// The Write method must be able to accept any amount
// of data, but it may operate more efficiently if all writes
// are a multiple of the block size.
// See hash.Hash interface.
func (h *Hash) BlockSize() int { return 1 }

// Reset resets the Hash to its initial state.
// See hash.Hash interface.
func (h *Hash) Reset() {
	h.curValue = h.table.InitCrc()
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
// See hash.Hash interface.
func (h *Hash) Sum(in []byte) []byte {
	s := h.CRC()
	for i := h.size; i > 0; {
		i--
		in = append(in, byte(s>>(8*i)))
	}
	return in
}

// Write implements io.Writer interface which is part of hash.Hash interface.
func (h *Hash) Write(p []byte) (n int, err error) {
	h.Update(p)
	return len(p), nil
}

// Update updates process supplied bytes and updates current (partial) CRC accordingly.
func (h *Hash) Update(p []byte) {
	h.curValue = h.table.UpdateCrc(h.curValue, p)
}

// CRC returns current CRC value for the data processed so far.
func (h *Hash) CRC() uint64 {
	return h.table.CRC(h.curValue)
}

// CalculateCRC is a convenience function allowing to calculate CRC in one call.
func (h *Hash) CalculateCRC(data []byte) uint64 {
	return h.table.CalculateCRC(data)
}

// NewHashWithTable creates a new Hash instance configured for table driven
// CRC calculation using a Table instance created elsewhere.
func NewHashWithTable(table *Table) *Hash {
	ret := &Hash{table: table}
	ret.size = (table.crcParams.Width + 7) / 8 // smalest number of bytes enough to store produced crc
	ret.Reset()
	return ret
}

// NewHash creates a new Hash instance configured for table driven
// CRC calculation according to parameters specified.
func NewHash(crcParams *Parameters) *Hash {
	return NewHashWithTable(NewTable(crcParams))
}

// CRC8 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (h *Hash) CRC8() uint8 {
	return h.table.CRC8(h.curValue)
}

// CRC16 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (h *Hash) CRC16() uint16 {
	return h.table.CRC16(h.curValue)
}

// CRC32 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (h *Hash) CRC32() uint32 {
	return h.table.CRC32(h.curValue)
}

// Table used by this Hash under the hood
func (h *Hash) Table() *Table {
	return h.table
}
