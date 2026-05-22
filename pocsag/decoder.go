package pocsag

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

// DecoderStats holds decoder performance counters.
type DecoderStats struct {
	SyncLosses      int64
	BCHFailures     int64
	MessagesDecoded int64
	DroppedSamples  int64
}

// Decoder converts IQ samples to AlphaMessages through the POCSAG pipeline.
type Decoder struct {
	stats    DecoderStats
	lastLock time.Time
}

// NewDecoder creates a new Decoder.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode processes a channel of IQ samples and outputs AlphaMessages.
func (d *Decoder) Decode(input <-chan complex64) <-chan obj.AlphaMessage {
	out := make(chan obj.AlphaMessage, 64)

	go func() {
		defer close(out)

		demod := newFMDemodulator()
		clock := newClockRecovery(22050, 512)
		sync := newFrameSynchronizer()
		assembler := newMessageAssembler()

		lockTimeout := 30 * time.Second
		lastPreambleCheck := time.Now()

		for sample := range input {
			// Stage 1: FM demodulation
			demod.feed(sample)
			demodVal := demod.output()

			// Stage 2: Clock recovery
			bit, haveBit := clock.feed(demodVal)
			if !haveBit {
				continue
			}

			// Check preamble lock
			if !clock.locked() {
				if time.Since(lastPreambleCheck) > lockTimeout {
					log.Printf("pocsag: clock recovery not locked for %v", lockTimeout)
					lastPreambleCheck = time.Now()
				}
				continue
			}

			// Stage 3: Frame sync -> codewords
			codewords := sync.processBits([]byte{bit})
			if len(codewords) == 0 {
				continue
			}

			d.lastLock = time.Now()

			// Stage 4-5: BCH decode + Message assembly
			for _, cw := range codewords {
				if cw == idleCodeword {
					continue
				}

				// Decode for error tracking (feedCodeword decodes internally)
				_, err := bchDecode(cw)
				if err != nil {
					atomic.AddInt64(&d.stats.BCHFailures, 1)
					failureRate := float64(d.stats.BCHFailures) / float64(d.stats.MessagesDecoded+1)
					if failureRate > 0.1 {
						log.Printf("pocsag: BCH failure rate %.1f%%", failureRate*100)
					}
					continue
				}

				msg := assembler.feedCodeword(cw)
				if msg != nil && msg.Valid {
					atomic.AddInt64(&d.stats.MessagesDecoded, 1)
					out <- *msg
				}
			}
		}
	}()

	return out
}

// Stats returns a snapshot of decoder performance counters.
func (d *Decoder) Stats() DecoderStats {
	return DecoderStats{
		SyncLosses:      atomic.LoadInt64(&d.stats.SyncLosses),
		BCHFailures:     atomic.LoadInt64(&d.stats.BCHFailures),
		MessagesDecoded: atomic.LoadInt64(&d.stats.MessagesDecoded),
		DroppedSamples:  atomic.LoadInt64(&d.stats.DroppedSamples),
	}
}
