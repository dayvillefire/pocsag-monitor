package pocsag

import "testing"

// BCH(31,21) generator polynomial: x^10 + x^9 + x^8 + x^6 + x^5 + x^3 + 1
// Binary: 11101101001 = 0x769

func Test_BCH_NoErrors(t *testing.T) {
	data, err := bchDecode(0x00000000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != 0 {
		t.Errorf("expected 0, got %d", data)
	}
}

func Test_BCH_EncodeDecode(t *testing.T) {
	cw := bchEncode(0x123456)
	data, err := bchDecode(cw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != 0x123456 {
		t.Errorf("expected 0x123456, got 0x%x", data)
	}
}

func Test_BCH_SingleBitError(t *testing.T) {
	cw := bchEncode(0x0AAAAA)
	corrupted := cw ^ (1 << 5)
	data, err := bchDecode(corrupted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != 0x0AAAAA {
		t.Errorf("single-bit error not corrected: expected 0x0AAAAA, got 0x%x", data)
	}
}

func Test_BCH_TwoBitErrors(t *testing.T) {
	cw := bchEncode(0x155555)
	corrupted := cw ^ (1 << 3) ^ (1 << 27)
	data, err := bchDecode(corrupted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != 0x155555 {
		t.Errorf("two-bit error not corrected: expected 0x155555, got 0x%x", data)
	}
}

func Test_BCH_Uncorrectable(t *testing.T) {
	cw := bchEncode(0x1FFFFF)
	corrupted := cw ^ (1 << 0) ^ (1 << 1) ^ (1 << 2)
	_, err := bchDecode(corrupted)
	if err == nil {
		t.Error("expected error for >2 bit errors")
	}
}

func Test_FrameSync_FindsSyncCodeword(t *testing.T) {
	preamble := generate1010Bits(576)
	syncWord := uint32(0x7CD215D8)
	bits := append(preamble, uint32ToBits(syncWord)...)
	bits = append(bits, uint32ToBits(0x7A89C197)...)
	bits = append(bits, uint32ToBits(0x7A89C197)...)
	fs := newFrameSynchronizer()
	results := fs.processBits(bits)
	if len(results) != 2 {
		t.Fatalf("expected 2 codewords, got %d", len(results))
	}
}

func Test_FrameSync_PartialErrorsInSync(t *testing.T) {
	preamble := generate1010Bits(576)
	corruptedSync := uint32(0x7CD215D8) ^ (1 << 3) ^ (1 << 17)
	bits := append(preamble, uint32ToBits(corruptedSync)...)
	bits = append(bits, uint32ToBits(0x7A89C197)...)
	fs := newFrameSynchronizer()
	results := fs.processBits(bits)
	if len(results) != 1 {
		t.Fatalf("expected 1 codeword with 2-bit sync errors, got %d", len(results))
	}
}
