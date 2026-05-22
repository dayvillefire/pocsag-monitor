# POCSAG-MONITOR

Golang POCSAG512 monitor, designed to be run on a Linux machine.

This code has been run on a [Raspberry Pi Zero W](https://amzn.to/3bflIyP) (the slower first-generation model) if anyone wants to run this on a standalone single-chip machine.

POCSAG monitor is now designed to be used with [pocsag-router](https://github.com/dayvillefire/pocsag-router), which allows aggregation and deduplication from multiple redundant monitors.

## Architecture

The monitor uses a **native Go pipeline** with no external binary dependencies:

```
RTL-SDR Dongle → sdr package (USB) → pocsag package (POCSAG512 decoder) → AlphaMessage
                                                    ↓
                                          router.Publish (NATS) or stdout
```

### Packages

- **`sdr/`** — Pure Go RTL2832U + R820T2 USB driver. Interfaces with the SDR hardware via Linux USBFS, no CGo or librtlsdr required. Also includes `FileIQSource` for replay testing with recorded IQ captures.
- **`pocsag/`** — Native POCSAG512 decoder pipeline: FM demodulator (arctan discriminator), low-pass filter, preamble-based clock recovery, frame synchronization, BCH(31,21) error correction, and message assembly into `AlphaMessage` structs.
- **`obj/`** — Shared types: `AlphaMessage` struct and text-format parser.
- **`config/`** — YAML configuration with defaults.

### Commands

| Command | Description |
|---------|-------------|
| `cmd/pocsag-monitor/` | Full service with web API, NATS router integration, and config file support |
| `cmd/pocsag-rx/` | Minimal CLI: streams decoded POCSAG messages to stdout |

### Pipeline Stages

```
IQ samples → FM Demod → LPF → Clock Recovery → Frame Sync → BCH(31,21) Decode → Message Assembly
(complex64)  (float64)         (bits)           (codewords)   (data words)       (AlphaMessage)
```

- **FM Demodulator:** Arctan-based discriminator measuring instantaneous frequency deviation
- **Low-Pass Filter:** Single-pole IIR, cutoff ~600 Hz (above POCSAG512's 256 Hz Nyquist)
- **Clock Recovery:** Locks to the 576-bit alternating preamble (`101010...`) at 512 bps
- **Frame Sync:** Detects sync codeword `0x7CD215D8` with up to 2-bit error tolerance
- **BCH(31,21):** Generator polynomial `x^10 + x^9 + x^8 + x^6 + x^5 + x^3 + 1`, corrects up to 2 bit errors per 31-bit codeword
- **Message Assembly:** Extracts capcode (18-bit) and function bits from address codewords, accumulates 7-bit ASCII text from message codewords

## pocsag-rx (Simple CLI)

Minimal stdout monitor. No config file, no web server, no NATS.

```bash
# Build
go build ./cmd/pocsag-rx/

# Run (requires SDR hardware)
./pocsag-rx -f 152.00750M

# With options
./pocsag-rx -f 155.16000M -p 2 -g 20
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `-f` | `152.00750M` | Frequency (MHz) |
| `-p` | `0` | PPM frequency correction |
| `-g` | `0` | Tuner gain (0 = auto) |
| `-r` | `22050` | Sample rate |
| `-d` | `0` | SDR device index |

**Output format:** `CAP: <capcode>  MSG: <message text>`

## pocsag-monitor (Full Service)

Full-featured service with web API and NATS integration.

### Quick Start

```bash
# Build
go build ./cmd/pocsag-monitor/

# Copy and edit config
cp config.yaml.example config.yaml
# Edit config.yaml with your settings

# Run
./pocsag-monitor
```

### Configuration

```yaml
debug: false
frequency: "152.00750M"
ppm: 0
api-port: 8001
instance-name: "MONITOR1"
sdr:
  device-index: 0
  gain: 0          # 0 = auto
  sample-rate: 22050
router:
  url: "tls://router-host:4222"
  topic: "pages"
```

### API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /api/debug` | Runtime info, environment, decoder stats |
| `GET /api/config` | Current configuration |
| `GET /api/version` | Build version |
| `GET /api/test/page/:capcode/:msg` | Inject test page |

### Systemd Service

A systemd service file is included at `pocsag-monitor.service`. Edit paths and user as needed.

```bash
sudo cp pocsag-monitor.service /etc/systemd/system/
sudo systemctl enable pocsag-monitor
sudo systemctl start pocsag-monitor
```

## Build

```bash
# Native (x86-64)
go build ./...

# Cross-compile for Raspberry Pi (ARMv6, no CGo)
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build ./cmd/pocsag-monitor/
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build ./cmd/pocsag-rx/
```

The build is CGo-free — produces statically linked binaries with no C library dependencies.

## Testing

```bash
# Run all tests
go test github.com/dayvillefire/pocsag-monitor/obj \
         github.com/dayvillefire/pocsag-monitor/pocsag \
         github.com/dayvillefire/pocsag-monitor/config

# POCSAG decoder tests (BCH, frame sync, demod, clock recovery, message assembly)
go test github.com/dayvillefire/pocsag-monitor/pocsag -v

# Recorded IQ replay test (requires capture file)
# rtl_fm -f 152.00750M -s 22050 > testdata/capture.iq
```

## PREREQUISITES

* Linux-compatible SDR (RTL2832U-based, e.g. NESDR Mini)
* Linux with USBFS support (`/dev/bus/usb`)

No external binaries (`rtl_fm`, `multimon-ng`, `librtlsdr`) are required — the monitor interfaces with the SDR hardware directly.

## RECOMMENDED HARDWARE

* [Raspberry Pi Zero W](https://amzn.to/3bflIyP) - 79$ US
* [NESDR Mini SDR](https://amzn.to/3TXecta) - 26$ US (the antenna that comes with this unit is not great, use the replacement one)
* [VHF Antenna](https://amzn.to/3ssavjt) - 9$ US
* [High Endurance SD Card](https://amzn.to/3Szn8Uj) - 10$ US

You'll also need a power supply. If the Pi Zero W is a little too expensive for you due to the parts shortage, consider "Le Potato", instead of the Pi Zero W:

* [Le Potato](https://amzn.to/3N57YF4) - 45$ US
* [Edimax Wifi Dongle](https://amzn.to/3SDyiaF) - 9$ US
* [Case and Power Supply](https://amzn.to/3gKwPlP) - 15$ US
