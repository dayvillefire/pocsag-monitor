package pocsag

import "math"

// fmDemodulator performs arctan-based FSK demodulation on IQ samples.
type fmDemodulator struct {
	prevI  float32
	prevQ  float32
	prevIn bool
	lpf    *lowPassFilter
}

func newFMDemodulator() *fmDemodulator {
	return &fmDemodulator{
		lpf: newLowPassFilter(0.05),
	}
}

func (d *fmDemodulator) feed(sample complex64) {
	i := real(sample)
	q := imag(sample)
	if !d.prevIn {
		d.prevI = i
		d.prevQ = q
		d.prevIn = true
		return
	}
	dot := i*d.prevI + q*d.prevQ
	cross := q*d.prevI - i*d.prevQ
	phase := math.Atan2(float64(cross), float64(dot))
	d.lpf.feed(phase)
	d.prevI = i
	d.prevQ = q
}

func (d *fmDemodulator) output() float64 {
	return d.lpf.output()
}

// lowPassFilter is a single-pole IIR low-pass filter.
// y[n] = alpha*x[n] + (1-alpha)*y[n-1]
type lowPassFilter struct {
	alpha float64
	value float64
}

func newLowPassFilter(alpha float64) *lowPassFilter {
	return &lowPassFilter{alpha: alpha}
}

func (f *lowPassFilter) feed(x float64) {
	f.value = f.alpha*x + (1-f.alpha)*f.value
}

func (f *lowPassFilter) output() float64 {
	return f.value
}
