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

		// Test helper methods return correct values as well
		if crcParams.Width == 8 {
			crc8 := tableDriven.CRC8()
			if crc8 != uint8(crc&0x00FF) {
				t.Errorf("Incorrect CRC8 0x%02x retrived %s (should be 0x%02x)", crc8, data, crc)
			}
		} else if crcParams.Width == 16 {
			crc16 := tableDriven.CRC16()
			if crc16 != uint16(crc&0x00FFFF) {
				t.Errorf("Incorrect CRC16 0x%04x retrived %s (should be 0x%04x)", crc16, data, crc)
			}
		} else if crcParams.Width == 32 {
			crc32 := tableDriven.CRC32()
			if crc32 != uint32(crc&0x00FFFFFFFF) {
				t.Errorf("Incorrect CRC8 0x%08x retrived %s (should be 0x%08x)", crc32, data, crc)
			}
		}

		// Test Hash's table directly and see there is no difference
		table := tableDriven.Table()
		calculated = table.CalculateCRC([]byte(data))
		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}

	}

	doTest(X25, "123456789", 0x906E)
	doTest(X25, "12345678901234567890", 0xA286)
	doTest(X25, "Introduction on CRC calculations", 0xF9B6)
	doTest(X25, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0x68B1)

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

	// More tests for various CRC algorithms (copied from java version)
	longText := "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits."

	testArrayData := make([]byte, 256)
	for i := 0; i < len(testArrayData); i++ {
		testArrayData[i] = byte(i & 0x0FF)
	}
	testArray := string(testArrayData)
	if len(testArray) != 256 {
		t.Fatalf("Logic error")
	}

	// merely a helper to make copying Spock test sets from java version of this library a bit easier
	doTestWithParameters := func(width uint, polynomial uint64, init uint64, reflectIn bool, reflectOut bool, finalXor uint64, crc uint64, testData string) {
		doTest(&Parameters{Width: width, Polynomial: polynomial, Init: init, ReflectIn: reflectIn, ReflectOut: reflectOut, FinalXor: finalXor}, testData, crc)
	}

	doTestWithParameters(3, 0x03, 0x00, false, false, 0x7, 0x04, "123456789") // CRC-3/GSM
	doTestWithParameters(3, 0x03, 0x00, false, false, 0x7, 0x06, longText)
	doTestWithParameters(3, 0x03, 0x00, false, false, 0x7, 0x02, testArray)
	doTestWithParameters(3, 0x03, 0x07, true, true, 0x0, 0x06, "123456789") // CRC-3/ROHC
	doTestWithParameters(3, 0x03, 0x07, true, true, 0x0, 0x03, longText)
	doTestWithParameters(4, 0x03, 0x00, true, true, 0x0, 0x07, "123456789")   // CRC-4/ITU
	doTestWithParameters(4, 0x03, 0x0f, false, false, 0xf, 0x0b, "123456789") // CRC-4/INTERLAKEN
	doTestWithParameters(4, 0x03, 0x0f, false, false, 0xf, 0x01, longText)    // CRC-4/INTERLAKEN
	doTestWithParameters(4, 0x03, 0x0f, false, false, 0xf, 0x07, testArray)   // CRC-4/INTERLAKEN
	doTestWithParameters(5, 0x09, 0x09, false, false, 0x0, 0x00, "123456789") // CRC-5/EPC
	doTestWithParameters(5, 0x15, 0x00, true, true, 0x0, 0x07, "123456789")   // CRC-5/ITU
	doTestWithParameters(6, 0x27, 0x3f, false, false, 0x0, 0x0d, "123456789") // CRC-6/CDMA2000-A
	doTestWithParameters(6, 0x07, 0x3f, false, false, 0x0, 0x3b, "123456789") // CRC-6/CDMA2000-B
	doTestWithParameters(6, 0x07, 0x3f, false, false, 0x0, 0x24, testArray)   // CRC-6/CDMA2000-B
	doTestWithParameters(7, 0x09, 0x00, false, false, 0x0, 0x75, "123456789") // CRC-7
	doTestWithParameters(7, 0x09, 0x00, false, false, 0x0, 0x78, testArray)   // CRC-7
	doTestWithParameters(7, 0x4f, 0x7f, true, true, 0x0, 0x53, "123456789")   // CRC-7/ROHC

	doTestWithParameters(8, 0x07, 0x00, false, false, 0x00, 0xf4, "123456789") // CRC-8
	doTestWithParameters(8, 0xa7, 0x00, true, true, 0x00, 0x26, "123456789")   // CRC-8/BLUETOOTH
	doTestWithParameters(8, 0x07, 0x00, false, false, 0x55, 0xa1, "123456789") // CRC-8/ITU
	doTestWithParameters(8, 0x9b, 0x00, true, true, 0x00, 0x25, "123456789")   // CRC-8/WCDMA
	doTestWithParameters(8, 0x31, 0x00, true, true, 0x00, 0xa1, "123456789")   // CRC-8/MAXIM

	doTestWithParameters(10, 0x233, 0x000, false, false, 0x000, 0x199, "123456789") // CRC-10

	doTestWithParameters(12, 0xd31, 0x00, false, false, 0xfff, 0x0b34, "123456789")   // CRC-12/GSM
	doTestWithParameters(12, 0x80f, 0x00, false, true, 0x00, 0x0daf, "123456789")     // CRC-12/UMTS
	doTestWithParameters(13, 0x1cf5, 0x00, false, false, 0x00, 0x04fa, "123456789")   // CRC-13/BBC
	doTestWithParameters(14, 0x0805, 0x00, true, true, 0x00, 0x082d, "123456789")     // CRC-14/DARC
	doTestWithParameters(14, 0x202d, 0x00, false, false, 0x3fff, 0x30ae, "123456789") // CRC-14/GSM

	doTestWithParameters(15, 0x4599, 0x00, false, false, 0x00, 0x059e, "123456789") // CRC-15
	doTestWithParameters(15, 0x4599, 0x00, false, false, 0x00, 0x2857, longText)
	doTestWithParameters(15, 0x6815, 0x00, false, false, 0x0001, 0x2566, "123456789") // CRC-15/MPT1327

	doTestWithParameters(21, 0x102899, 0x000000, false, false, 0x000000, 0x0ed841, "123456789") // CRC-21/CAN-FD
	doTestWithParameters(24, 0x864cfb, 0xb704ce, false, false, 0x000000, 0x21cf02, "123456789") // CRC-24
	doTestWithParameters(24, 0x5d6dcb, 0xfedcba, false, false, 0x000000, 0x7979bd, "123456789") // CRC-24/FLEXRAY-A
	doTestWithParameters(24, 0x00065b, 0x555555, true, true, 0x000000, 0xc25a56, "123456789")   // "CRC-24/BLE"

	doTestWithParameters(31, 0x04c11db7, 0x7fffffff, false, false, 0x7fffffff, 0x0ce9e46c, "123456789") // CRC-31/PHILIPS
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
