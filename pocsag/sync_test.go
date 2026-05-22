package pocsag

import "testing"

func Test_ClockRecovery_LocksToPreamble(t *testing.T) {
	cr := newClockRecovery(22050, 512)
	samplesPerBit := 22050.0 / 512.0
	for i := 0; i < 600*int(samplesPerBit); i++ {
		bitIdx := i / int(samplesPerBit)
		var val float64
		if bitIdx%2 == 0 {
			val = 1.0
		} else {
			val = -1.0
		}
		cr.feed(val)
	}
	if !cr.locked() {
		t.Error("clock recovery should lock to 1010... preamble")
	}
}

func Test_ClockRecovery_ProducesBits(t *testing.T) {
	cr := newClockRecovery(22050, 512)
	samplesPerBit := 22050.0 / 512.0
	bits := make([]byte, 0)
	for i := 0; i < 700*int(samplesPerBit); i++ {
		bitIdx := i / int(samplesPerBit)
		var val float64
		if bitIdx < 576 {
			if bitIdx%2 == 0 {
				val = 1.0
			} else {
				val = -1.0
			}
		} else {
			if bitIdx == 576 {
				val = -1.0
			} else {
				val = 1.0
			}
		}
		if b, ok := cr.feed(val); ok {
			bits = append(bits, b)
		}
	}
	if len(bits) == 0 {
		t.Fatal("clock recovery should produce bits after locking")
	}
	t.Logf("got %d bits total", len(bits))
}
