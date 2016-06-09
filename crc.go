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
func CalculateCRC(crcParams *Parameters, data []byte) uint64 {

	curValue := crcParams.Init
	topbit := uint64(1) << (crcParams.Width - 1)
	mask := (topbit << 1) - 1

	for i := 0; i < len(data); i++ {
		var curByte = uint64(data[i]) & 0x00FF
		if crcParams.ReflectIn {
			curByte = reflect(curByte, 8)
		}
		curValue ^= (curByte << (crcParams.Width - 8))
		for j := 0; j < 8; j++ {
			if (curValue & topbit) != 0 {
				curValue = (curValue << 1) ^ crcParams.Polynomial
			} else {
				curValue = (curValue << 1)
			}
		}

	}
	if crcParams.ReflectOut {
		curValue = reflect(curValue, crcParams.Width)
	}

	curValue = curValue ^ crcParams.FinalXor

	return curValue & mask
}

// Hash represents the partial evaluation of a checksum using table-driven
// implementation. It also implements hash.Hash interface.
type Hash struct {
	crcParams Parameters
	crctable  []uint64
	curValue  uint64
	mask      uint64
	size      uint
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
	h.curValue = h.crcParams.Init
	if h.crcParams.ReflectIn {
		h.curValue = reflect(h.crcParams.Init, h.crcParams.Width)
	}
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
	if h.crcParams.ReflectIn {
		for _, v := range p {
			h.curValue = h.crctable[(byte(h.curValue)^v)&0xFF] ^ (h.curValue >> 8)
		}
	} else {
		for _, v := range p {
			h.curValue = h.crctable[((byte(h.curValue>>(h.crcParams.Width-8))^v)&0xFF)] ^ (h.curValue << 8)
		}
	}
}

// CRC returns current CRC value for the data processed so far.
func (h *Hash) CRC() uint64 {
	ret := h.curValue

	if h.crcParams.ReflectOut != h.crcParams.ReflectIn {
		ret = reflect(ret, h.crcParams.Width)
	}
	return (ret ^ h.crcParams.FinalXor) & h.mask
}

// CalculateCRC is a convenience function allowing to calculate CRC in one call.
func (h *Hash) CalculateCRC(data []byte) uint64 {
	h.Reset() // just in case
	h.Update(data)
	return h.CRC()
}

// NewHash creates a new Hash instance configured for table driven
// CRC calculation according to parameters specified.
func NewHash(crcParams *Parameters) *Hash {
	ret := &Hash{crcParams: *crcParams}
	ret.mask = (uint64(1) << crcParams.Width) - 1
	ret.size = (crcParams.Width + 7) / 8 // smalest number of bytes enough to store produced crc
	ret.crctable = make([]uint64, 256, 256)

	tmp := make([]byte, 1, 1)
	tableParams := *crcParams
	tableParams.Init = 0
	tableParams.ReflectOut = tableParams.ReflectIn
	tableParams.FinalXor = 0
	for i := 0; i < 256; i++ {
		tmp[0] = byte(i)
		ret.crctable[i] = CalculateCRC(&tableParams, tmp)
	}
	ret.Reset()

	return ret
}

// CRC16 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (h *Hash) CRC16() uint16 {
	return uint16(h.CRC())
}

// CRC32 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (h *Hash) CRC32() uint32 {
	return uint32(h.CRC())
}
