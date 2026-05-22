// sdr/source.go
package sdr

// IQSample is a single complex IQ sample.
type IQSample = complex64

// IQSource provides IQ samples from an SDR device or file.
type IQSource interface {
	// Start begins streaming IQ samples. freqHz is the center frequency in Hz.
	// sampleRate is samples per second. ppm is frequency correction.
	// gain is tuner gain (0 = auto).
	Start(freqHz int, sampleRate int, ppm int, gain int) (<-chan IQSample, error)
	// Stop halts streaming and releases resources.
	Stop() error
}
