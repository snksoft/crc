package crc

import (
	"hash"
	"testing"
)

func TestCRCAlgorithms(t *testing.T) {

	doTest := func(crcParams *Parameters, data string, crc uint64) {
		calculated := CalculateCRC(crcParams, []byte(data))
		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}

		// same test using table driven
		tableDriven := NewHash(crcParams)
		calculated = tableDriven.CalculateCRC([]byte(data))
		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}

		// same test feeding data in chunks of different size
		tableDriven.Reset()
		var start = 0
		var step = 1
		for start < len(data) {
			end := start + step
			if end > len(data) {
				end = len(data)
			}
			tableDriven.Update([]byte(data[start:end]))
			start = end
			step *= 2
		}
		calculated = tableDriven.CRC()
		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}
	}

	doTest(CCITT, "123456789", 0x29B1)
	doTest(CCITT, "12345678901234567890", 0xDA31)
	doTest(CCITT, "Introduction on CRC calculations", 0xC87E)
	doTest(CCITT, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0xD6ED)

	doTest(XMODEM, "123456789", 0x31C3)
	doTest(XMODEM, "12345678901234567890", 0x2C89)
	doTest(XMODEM, "Introduction on CRC calculations", 0x3932)
	doTest(XMODEM, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0x4E86)

	doTest(XMODEM2, "123456789", 0x0C73)
	doTest(XMODEM2, "12345678901234567890", 0x122E)
	doTest(XMODEM2, "Introduction on CRC calculations", 0x0638)
	doTest(XMODEM2, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0x187A)

	doTest(CRC32, "123456789", 0xCBF43926)
	doTest(CRC32, "12345678901234567890", 0x906319F2)
	doTest(CRC32, "Introduction on CRC calculations", 0x814F2B45)
	doTest(CRC32, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0x8F273817)

	doTest(Castagnoli, "123456789", 0xE3069283)
	doTest(Castagnoli, "12345678901234567890", 0xA8B4A6B9)
	doTest(Castagnoli, "Introduction on CRC calculations", 0x54F98A9E)
	doTest(Castagnoli, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0x864FDAFC)

	doTest(Koopman, "123456789", 0x2D3DD0AE)
	doTest(Koopman, "12345678901234567890", 0xCC53DEAC)
	doTest(Koopman, "Introduction on CRC calculations", 0x1B8101F9)
	doTest(Koopman, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0xA41634B2)

	doTest(CRC64ISO, "123456789", 0xB90956C775A41001)
	doTest(CRC64ISO, "12345678901234567890", 0x8DB93749FB37B446)
	doTest(CRC64ISO, "Introduction on CRC calculations", 0xBAA81A1ED1A9209B)
	doTest(CRC64ISO, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0x347969424A1A7628)

	doTest(CRC64ECMA, "123456789", 0x995DC9BBDF1939FA)
	doTest(CRC64ECMA, "12345678901234567890", 0x0DA1B82EF5085A4A)
	doTest(CRC64ECMA, "Introduction on CRC calculations", 0xCF8C40119AE90DCB)
	doTest(CRC64ECMA, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0x31610F76CFB272A5)
}

func TestSizeMethods(t *testing.T) {
	testWidth := func(width uint, expectedSize int) {
		h := NewHash(&Parameters{Width: width, Polynomial: 1})
		s := h.Size()
		if s != expectedSize {
			t.Errorf("Incorrect Size calculated for width %d:  %d when should be %d", width, s, expectedSize)
		}
		bs := h.BlockSize()
		if bs != 1 {
			t.Errorf("Incorrect Block Size returned for width %d:  %d when should always be 1", width, bs)
		}
	}

	testWidth(3, 1)
	testWidth(8, 1)
	testWidth(12, 2)
	testWidth(16, 2)
	testWidth(32, 4)
	testWidth(64, 8)

}

func TestHashInterface(t *testing.T) {
	doTest := func(crcParams *Parameters, data string, crc uint64) {
		// same test using table driven
		var h hash.Hash = NewHash(crcParams)

		// same test feeding data in chunks of different size
		h.Reset()
		var start = 0
		var step = 1
		for start < len(data) {
			end := start + step
			if end > len(data) {
				end = len(data)
			}
			h.Write([]byte(data[start:end]))
			start = end
			step *= 2
		}

		buf := make([]byte, 0, 0)
		buf = h.Sum(buf)

		if len(buf) != h.Size() {
			t.Errorf("Wrong number of bytes appended by Sum(): %d when should be %d", len(buf), h.Size())
		}

		calculated := uint64(0)
		for _, b := range buf {
			calculated <<= 8
			calculated += uint64(b)
		}

		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}
	}

	doTest(&Parameters{Width: 8, Polynomial: 0x07, Init: 0x00, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}, "123456789", 0xf4)
	doTest(CCITT, "12345678901234567890", 0xDA31)
	doTest(CRC64ECMA, "Introduction on CRC calculations", 0xCF8C40119AE90DCB)
	doTest(CRC32C, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0x864FDAFC)
}

func BenchmarkCCITT(b *testing.B) {
	data := []byte("Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.")
	for i := 0; i < b.N; i++ {
		tableDriven := NewHash(CCITT)
		tableDriven.Update(data)
		tableDriven.CRC()
	}
}
