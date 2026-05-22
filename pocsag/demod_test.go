package pocsag

import (
	"math"
	"testing"
)

func Test_FMDemod_PositiveDeviation(t *testing.T) {
	d := newFMDemodulator()
	phaseInc := 2 * math.Pi * 4500 / 22050
	for i := 0; i < 100; i++ {
		phase := float32(phaseInc * float64(i))
		sample := complex(float32(math.Cos(float64(phase))), float32(math.Sin(float64(phase))))
		d.feed(sample)
	}
	out := d.output()
	if out <= 0 {
		t.Errorf("expected positive output for positive deviation, got %f", out)
	}
}

func Test_FMDemod_NegativeDeviation(t *testing.T) {
	d := newFMDemodulator()
	phaseInc := 2 * math.Pi * -4500 / 22050
	for i := 0; i < 100; i++ {
		phase := float32(phaseInc * float64(i))
		sample := complex(float32(math.Cos(float64(phase))), float32(math.Sin(float64(phase))))
		d.feed(sample)
	}
	out := d.output()
	if out >= 0 {
		t.Errorf("expected negative output for negative deviation, got %f", out)
	}
}

func Test_FMDemod_NoModulation(t *testing.T) {
	d := newFMDemodulator()
	sample := complex(float32(1.0), float32(0.0))
	for i := 0; i < 100; i++ {
		d.feed(sample)
	}
	out := d.output()
	if math.Abs(float64(out)) > 0.01 {
		t.Errorf("expected near-zero output, got %f", out)
	}
}
