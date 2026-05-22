# Native POCSAG Pipeline Design

Replace `rtl_fm` and `multimon-ng` external C binaries with pure Go implementations, creating a single-binary pipeline from SDR IQ samples to decoded AlphaMessage output.

## Scope

- Pure Go RTL2832U + R820T2 SDR interface (no CGo, no librtlsdr)
- Native POCSAG512 decoder (FSK demod, bit sync, frame decode, BCH error correction, message assembly)
- Decoder outputs `obj.AlphaMessage` directly (no text intermediary)
- POCSAG512 only (not 1200 or 2400 bps)
- Both pieces must be complete — full replacement of both binaries

## Package Structure

```
pocsag-monitor/
├── cmd/pocsag-monitor/     # main.go, router.go, api.go (existing, modified)
├── config/                 # config.go (modified — remove binary fields, add SDR block)
├── obj/                    # AlphaMessage, ParseAlphaMessage (keep ParseAlphaMessage)
├── sdr/                    # NEW: RTL-SDR hardware interface
│   ├── source.go           # IQSource interface + IQSample type
│   ├── rtl2832.go          # RTL2832U + R820T2 USB implementation
│   └── file.go             # FileIQSource for replay testing
└── pocsag/                 # NEW: POCSAG512 decoder pipeline
    ├── decoder.go          # Decoder interface + orchestration + stats
    ├── demod.go            # FSK demodulator (arctan discriminator)
    ├── sync.go             # Bit synchronization / clock recovery
    ├── frame.go            # Frame detection + BCH(31,21) error correction
    ├── message.go          # Message reassembly → obj.AlphaMessage
    └── decoder_test.go     # Integration tests with IQ fixtures
```

## Key Interfaces

```go
// sdr/source.go
type IQSource interface {
    Start(freqHz int, sampleRate int, ppm int, gain int) (<-chan complex64, error)
    Stop() error
}

// pocsag/decoder.go
type Decoder interface {
    Decode(input <-chan complex64) <-chan obj.AlphaMessage
    Stats() DecoderStats
}
```

## Architecture

Interface-driven pipeline stages. Each processing stage is independently testable with known inputs/outputs.

```
IQSource → FSKDemodulator → LowPassFilter → ClockRecovery → FrameSync → BCHDecoder → MessageAssembler
(complex64)   (float64)        (float64)       (bits)         (codewords)  (data words)  (AlphaMessage)
```

## POCSAG Decoder Pipeline

### 1. FM Demodulator
Arctan-based discriminator: `atan2(Q[n]*I[n-1] - I[n]*Q[n-1], I[n]*I[n-1] + Q[n]*Q[n-1])`. For POCSAG512 at 22050 sps, ~43 samples per bit. Nominal deviation ±4.5 kHz.

### 2. Low-Pass Filter
Moving average or single-pole IIR, cutoff ~600 Hz (above 512 bps Nyquist of 256 Hz, to preserve data transitions while removing high-frequency noise).

### 3. Clock Recovery
During preamble (576 bits of alternating 101010...), measure transition timing to lock the bit clock. Once locked, use for bit decisions. Handle minor drift.

### 4. Frame Synchronization
Scan for 32-bit sync codeword `0x7CD215D8`. Allow up to 2 bit errors in detection. Lock onto batch structure: 1 sync + 16 codewords (8 frames × 2 codewords each).

### 5. BCH(31,21) Decode
Generator polynomial: `x^10 + x^9 + x^8 + x^6 + x^5 + x^3 + 1` (CCIR Rec. 584). Syndrome calculation → error locator → correct up to 2 bits. Output: 21-bit data word + address/message flag (bit 0).

### 6. Message Assembly
- Address codeword (bit 0 = 0): extract capcode (bits 1-18) and function bits (19-20)
- Message codeword (bit 0 = 1): accumulate 20-bit chunks into message buffer
- On next address or idle codeword: emit completed AlphaMessage
- Translate function bits + message data into AlphaMessage using existing logic

## SDR Reader

### Device Discovery
USBFS via `/dev/bus/usb/XXX/XXX` using `golang.org/x/sys/unix` ioctl. Find device by USB VID:PID `0x0bda:0x2832` or `0x0bda:0x2838`. Claim interface 0, configure bulk endpoint for IQ streaming.

### RTL2832U Configuration
Control transfers to initialize demodulator: reset, set sample rate register, configure ADC, enable IQ output mode. Sample rate formula: 28.8 MHz / divider. Register sequences ported from librtlsdr/Osmocom documentation.

### R820T2 Tuner Programming
I2C passthrough via RTL2832U registers. Set PLL frequency, LNA/mixer/VGA gain stages. Known init sequence of ~25 registers. PPM correction applied as frequency offset.

### IQ Sample Streaming
Read from bulk endpoint into ring buffer. Samples: interleaved uint8 I/Q pairs `[I0, Q0, I1, Q1, ...]`. Convert to `complex64` via channel. Expected throughput: ~44 KB/s at 22050 sps (trivially manageable).

## Configuration Changes

### Removed
- `RtlFmBinary` (was `rtlfm`, default `"rtl_fm"`)
- `MultiMonBinary` (was `multimon`, default `"multimon-ng"`)

### Added
```go
SDR struct {
    DeviceIndex int `yaml:"device-index" default:"0"`
    Gain        int `yaml:"gain" default:"0"`       // 0 = auto
    SampleRate  int `yaml:"sample-rate" default:"22050"`
} `yaml:"sdr"`
```

### Retained
- `Frequency` (string, e.g. `"152.00750M"`)
- `PPM` (int, frequency correction)

## main.go Integration

### Before (~30 lines)
- `exec.Command` for `rtl_fm` and `multimon-ng`
- Pipe setup between binaries
- Signal handling for child processes
- `bufio.Scanner` + `ParseAlphaMessage` in hot loop

### After (~6 lines)
```go
src := sdr.NewRTLSource(cfg.SDR.DeviceIndex)
iqCh, err := src.Start(freqHz, cfg.SDR.SampleRate, cfg.PPM, cfg.SDR.Gain)
dec := pocsag.NewDecoder()
for alpha := range dec.Decode(iqCh) {
    router.Publish(cfg.Router.Topic, alpha)
}
```

### obj/parse.go
`ParseAlphaMessage` is retained — still used by `/api/test/page` endpoint which manually constructs AlphaMessage from URL parameters.

## Error Handling

### Soft errors (log + skip, don't crash)
- BCH decode failures (>2 bit errors): skip codeword, log when error rate exceeds 10%
- Frame sync loss: discard partial batch, return to preamble search
- Clock recovery unlock: keep scanning for preamble, log if >30s without lock
- Sample buffer overrun: drop oldest samples, increment counter

### Hard errors (fatal)
- SDR device not found at startup: fatal
- USB communication failure: fatal

### Observability
```go
type DecoderStats struct {
    SyncLosses      int64
    BCHFailures     int64
    MessagesDecoded int64
    DroppedSamples  int64
}
```
Exposed via existing `GET /api/debug` endpoint.

## Testing Strategy

### Unit tests (no hardware)
- BCH(31,21): tables of known codewords → expected data. Test 0/1/2/>2 bit error cases.
- Frame sync: bit sequences with sync codeword at known positions, with and without errors.
- FM demod: synthetic IQ at known frequencies → verify deviation output.
- Message assembly: address/message/idle codeword sequences → correct AlphaMessage.

### Integration tests (recorded IQ)
- Record 30-60s IQ from `rtl_fm -f <freq> -s 22050 > testdata/capture.iq`
- `FileIQSource` replays capture through full decoder pipeline
- Output compared against known `multimon-ng` output from same capture
- These are the golden tests for end-to-end correctness

### Existing tests
- `obj/parse_test.go` — continues passing (ParseAlphaMessage unchanged)
- `router_test.go` — continues passing (routing unchanged)
- `config_test.go` — update for new SDR fields, remove old binary fields

## Verification

1. `go build ./...` — compiles without CGo
2. `go test ./...` — all unit + integration tests pass
3. Manual test with live SDR hardware against known pager frequencies
4. Cross-compile for linux/arm (Raspberry Pi) and verify single binary runs without librtlsdr installed
