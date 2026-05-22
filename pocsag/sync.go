package pocsag

type clockRecovery struct {
	sampleRate    int
	bitRate       int
	samplesPerBit float64
	sampleCount   int
	prevSample    float64
	lockedIn      bool
	preambleBits  int
	bitOutput     []byte
	bitPhase      float64
	phaseLocked   bool
}

func newClockRecovery(sampleRate, bitRate int) *clockRecovery {
	return &clockRecovery{
		sampleRate:    sampleRate,
		bitRate:       bitRate,
		samplesPerBit: float64(sampleRate) / float64(bitRate),
		bitOutput:     make([]byte, 0, 64),
	}
}

func (cr *clockRecovery) feed(demodValue float64) (byte, bool) {
	cr.sampleCount++

	// Use the first zero crossing to establish the bit clock phase.
	// For the alternating POCSAG preamble a zero crossing occurs at every
	// bit boundary; updating bitPhase on each one would prevent the sample
	// counter from ever reaching a full bit period.
	if !cr.phaseLocked && cr.sampleCount > 1 {
		if (cr.prevSample >= 0 && demodValue < 0) || (cr.prevSample < 0 && demodValue >= 0) {
			frac := cr.prevSample / (cr.prevSample - demodValue)
			cr.bitPhase = float64(cr.sampleCount-1) + frac
			cr.phaseLocked = true
		}
	}
	cr.prevSample = demodValue

	// Produce a bit every samplesPerBit after the established phase.
	if cr.phaseLocked {
		bitPeriod := cr.samplesPerBit
		if float64(cr.sampleCount)-cr.bitPhase >= bitPeriod {
			var bit byte
			if demodValue >= 0 {
				bit = 1
			} else {
				bit = 0
			}
			cr.bitOutput = append(cr.bitOutput, bit)
			cr.bitPhase += bitPeriod
		}
	}

	// Before locking, accumulate bits and look for the 1010... preamble pattern.
	// Bits are NOT drained yet so the buffer can grow past the detection threshold.
	if !cr.lockedIn && len(cr.bitOutput) >= 32 {
		if cr.isPreamblePattern() {
			cr.preambleBits++
			if cr.preambleBits >= 32 {
				cr.lockedIn = true
			}
		} else {
			cr.preambleBits = 0
			cr.bitOutput = cr.bitOutput[:0]
		}
	}

	// After locking, drain one accumulated bit per call.
	if cr.lockedIn && len(cr.bitOutput) > 0 {
		b := cr.bitOutput[0]
		cr.bitOutput = cr.bitOutput[1:]
		return b, true
	}
	return 0, false
}

func (cr *clockRecovery) locked() bool {
	return cr.lockedIn
}

// isPreamblePattern checks whether the last 16 accumulated bits are
// alternating (1010... or 0101...). Unlike a phase-sensitive check, this
// accepts either starting polarity, which is correct for POCSAG's 576-bit
// alternating preamble.
func (cr *clockRecovery) isPreamblePattern() bool {
	n := len(cr.bitOutput)
	if n < 16 {
		return false
	}
	start := n - 16
	for i := start; i < n-1; i++ {
		if cr.bitOutput[i] == cr.bitOutput[i+1] {
			return false
		}
	}
	return true
}
