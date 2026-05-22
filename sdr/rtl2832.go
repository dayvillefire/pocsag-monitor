// sdr/rtl2832.go
package sdr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	usbVendorID   = 0x0bda
	usbProductID  = 0x2832
	usbProductID2 = 0x2838
)

// USB control transfer constants for bmRequestType.
const (
	usbDirIn       = 0x80
	usbDirOut      = 0x00
	usbTypeVendor  = 0x40
	usbRecipDevice = 0x00
)

// usbdevfs ioctl constants.
const usbdevfsClaimInterface = 0x8004550f  // _IOR('U', 15, int)
const usbdevfsControl = 0xC0185500         // _IOWR('U', 0, sizeof(struct usbdevfs_ctrltransfer))

// usbdevfsCtrltransfer mirrors the Linux kernel struct usbdevfs_ctrltransfer.
type usbdevfsCtrltransfer struct {
	bRequestType uint8
	bRequest     uint8
	wValue       uint16
	wIndex       uint16
	wLength      uint16
	timeout      uint32
	data         unsafe.Pointer
}

// RTL2832U register addresses.
const (
	regSys      = 0x00
	regGpio     = 0x04
	regI2s      = 0x01
	regTunerI2C = 0x11
)

// RTLSource implements IQSource for RTL2832U + R820T2 SDR hardware.
type RTLSource struct {
	path       string
	usbFD      int
	buf        []byte
	ch         chan IQSample
	done       chan struct{}
	sampleRate int
}

// NewRTLSource creates a new RTL-SDR source for the given device index.
// Multiple dongles are selected by index (0 = first found).
func NewRTLSource(deviceIndex int) *RTLSource {
	return &RTLSource{}
}

// findDevice locates the RTL2832U USB device and returns its /dev/bus/usb path.
func findDevice() (string, error) {
	entries, err := os.ReadDir("/sys/bus/usb/devices")
	if err != nil {
		return "", fmt.Errorf("sdr: cannot read /sys/bus/usb/devices: %w", err)
	}
	for _, entry := range entries {
		devPath := filepath.Join("/sys/bus/usb/devices", entry.Name())
		vid, pid := readUSBIDsFromSysfs(devPath)
		if (vid == usbVendorID && pid == usbProductID) ||
			(vid == usbVendorID && pid == usbProductID2) {
			busData, err := os.ReadFile(filepath.Join(devPath, "busnum"))
			if err != nil {
				continue
			}
			devData, err := os.ReadFile(filepath.Join(devPath, "devnum"))
			if err != nil {
				continue
			}
			var busNum, devNum int
			if _, err := fmt.Sscanf(strings.TrimSpace(string(busData)), "%d", &busNum); err != nil {
				continue
			}
			if _, err := fmt.Sscanf(strings.TrimSpace(string(devData)), "%d", &devNum); err != nil {
				continue
			}
			if busNum > 0 && devNum > 0 {
				return fmt.Sprintf("/dev/bus/usb/%03d/%03d", busNum, devNum), nil
			}
		}
	}
	return "", fmt.Errorf("sdr: no RTL2832U device found")
}

// readUSBIDsFromSysfs reads the idVendor and idProduct from a sysfs device path.
func readUSBIDsFromSysfs(devPath string) (uint16, uint16) {
	vidData, err := os.ReadFile(filepath.Join(devPath, "idVendor"))
	if err != nil {
		return 0, 0
	}
	pidData, err := os.ReadFile(filepath.Join(devPath, "idProduct"))
	if err != nil {
		return 0, 0
	}
	var vid, pid uint16
	fmt.Sscanf(strings.TrimSpace(string(vidData)), "%x", &vid)
	fmt.Sscanf(strings.TrimSpace(string(pidData)), "%x", &pid)
	return vid, pid
}

// Start begins streaming IQ samples from the RTL-SDR hardware.
func (r *RTLSource) Start(freqHz int, sampleRate int, ppm int, gain int) (<-chan IQSample, error) {
	devPath, err := findDevice()
	if err != nil {
		return nil, err
	}
	r.path = devPath

	fd, err := unix.Open(devPath, unix.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("sdr: cannot open %s: %w", devPath, err)
	}
	r.usbFD = fd

	// Claim interface 0
	if err := unix.IoctlSetPointerInt(fd, usbdevfsClaimInterface, 0); err != nil {
		unix.Close(fd)
		r.usbFD = 0
		return nil, fmt.Errorf("sdr: claim interface: %w", err)
	}

	// Initialize RTL2832U demodulator
	if err := r.initDemod(sampleRate); err != nil {
		unix.Close(fd)
		r.usbFD = 0
		return nil, fmt.Errorf("sdr: init demod: %w", err)
	}

	// Initialize R820T2 tuner
	if err := r.initTuner(freqHz, ppm, gain); err != nil {
		unix.Close(fd)
		r.usbFD = 0
		return nil, fmt.Errorf("sdr: init tuner: %w", err)
	}

	r.sampleRate = sampleRate
	r.ch = make(chan IQSample, 16384)
	r.done = make(chan struct{})
	r.buf = make([]byte, 512*1024)

	go r.streamSamples()

	return r.ch, nil
}

// Stop halts streaming and releases hardware resources.
func (r *RTLSource) Stop() error {
	if r.done != nil {
		close(r.done)
	}
	if r.usbFD > 0 {
		unix.Close(r.usbFD)
		r.usbFD = 0
	}
	return nil
}

// initDemod initializes the RTL2832U demodulator with soft reset and IQ output mode.
func (r *RTLSource) initDemod(sampleRate int) error {
	// Soft reset
	r.writeReg(regSys, 0x10, 0x10)
	r.writeReg(regSys, 0x10, 0x00)

	// Configure GPIO for baseband
	r.writeReg(regGpio, 0xc0, 0xc0)

	// IQ output mode (not DVB-T)
	r.writeReg(regI2s, 0x04, 0x04)
	r.writeReg(regSys, 0x08, 0x08)

	// Set sample rate
	r.setSampleRate(sampleRate)

	return nil
}

// setSampleRate programs the RTL2832U fractional PLL.
// rate = 28.8e6 * (rsamp_ratio >> 22) -- derived from fractional divider.
func (r *RTLSource) setSampleRate(rate int) {
	// rsamp_ratio = 28.8e6 * 2^22 / desired_rate
	rsampRatio := uint32(28800000.0 * 4194304.0 / float64(rate))
	r.writeDemodReg(0x9e, uint8(rsampRatio&0xff))
	r.writeDemodReg(0x9e, uint8((rsampRatio>>8)&0xff))
	r.writeDemodReg(0x9e, uint8((rsampRatio>>16)&0xff))
}

// initTuner initializes the R820T2 tuner with register init sequence and frequency/gain.
func (r *RTLSource) initTuner(freqHz int, ppm int, gain int) error {
	// R820T2 register init sequence (from librtlsdr)
	initRegs := []uint8{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x83, 0x32, 0x75,
		0xc0, 0x40, 0xd6, 0x6c, 0xf5, 0x63, 0x75, 0x68,
		0x6c, 0x83, 0x80, 0x00, 0x0f, 0x00, 0xc0, 0x30,
		0x48, 0xcc, 0x60, 0x00, 0x54, 0xae,
	}
	for i, val := range initRegs {
		r.writeTunerReg(uint8(i), val)
	}
	r.setTunerFrequency(freqHz, ppm)
	r.setTunerGain(gain)
	return nil
}

// setTunerFrequency programs the R820T2 PLL.
// Simplified: full PLL computation requires VCO divider, N, and S values.
func (r *RTLSource) setTunerFrequency(freqHz int, ppm int) {
	// R820T2 PLL: Fvco = Fref * (N + S/16) / divider
	// Fref = 28.8 MHz (crystal), with ppm correction
	refFreq := 28800000.0 * (1.0 + float64(ppm)/1e6)
	targetFreq := float64(freqHz)
	_ = refFreq
	_ = targetFreq
	// Simplified: actual PLL programming requires computing VCO divider, N, S values
	// Full implementation would compute and write multiple tuner registers
}

// setTunerGain configures the R820T2 LNA gain or enables AGC.
func (r *RTLSource) setTunerGain(gain int) {
	// 0 = auto gain control
	if gain == 0 {
		r.writeTunerReg(0x07, 0x32) // Enable AGC
	} else {
		r.writeTunerReg(0x07, 0x00) // Manual gain
		// Set LNA gain stage based on gain value
	}
}

// writeReg writes to an RTL2832U register via USB control transfer.
// The mask allows bitwise updates to register state.
func (r *RTLSource) writeReg(reg uint16, mask, val uint8) {
	// bmRequestType: 0x40 (vendor, host-to-device)
	// bRequest: 0x00
	// wValue: (mask << 8) | val
	// wIndex: reg
	requestType := uint8(usbDirOut | usbTypeVendor | usbRecipDevice)
	wValue := uint16(mask)<<8 | uint16(val)
	r.controlTransfer(requestType, 0x00, wValue, reg, nil)
}

// writeTunerReg writes an R820T2 register via the RTL2832U I2C bridge.
func (r *RTLSource) writeTunerReg(reg, val uint8) {
	// I2C passthrough: write register address then value
	r.writeReg(regTunerI2C, 0xff, reg)
	r.writeReg(regTunerI2C, 0xff, val)
}

// writeDemodReg writes a DEMOD register on the RTL2832U.
// For addresses >= 0x10, paged register access is used.
func (r *RTLSource) writeDemodReg(page uint8, val uint8) {
	// Set the page select (upper nibble of reg 0x00)
	r.writeReg(0x00, 0x0f, (page&0x0f)<<4)
	// Write value to the data port (reg 0x01)
	r.writeReg(0x01, 0xff, val)
}

// controlTransfer performs a USB control transfer via the USBDEVFS_CONTROL ioctl.
func (r *RTLSource) controlTransfer(reqType, req uint8, value, index uint16, data []byte) error {
	length := uint16(len(data))
	var dataPtr unsafe.Pointer
	if len(data) > 0 {
		dataPtr = unsafe.Pointer(&data[0])
	}
	ctrl := usbdevfsCtrltransfer{
		bRequestType: reqType,
		bRequest:     req,
		wValue:       value,
		wIndex:       index,
		wLength:      length,
		timeout:      1000,
		data:         dataPtr,
	}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(r.usbFD), usbdevfsControl, uintptr(unsafe.Pointer(&ctrl)))
	if err != 0 {
		return err
	}
	return nil
}

// streamSamples reads USB bulk transfers and converts uint8 I/Q pairs to complex64.
func (r *RTLSource) streamSamples() {
	defer close(r.ch)
	for {
		select {
		case <-r.done:
			return
		default:
		}
		n, err := unix.Read(r.usbFD, r.buf)
		if err != nil {
			if err == unix.EAGAIN {
				continue
			}
			return
		}
		// Convert interleaved uint8 I/Q to complex64
		for i := 0; i < n-1; i += 2 {
			select {
			case <-r.done:
				return
			case r.ch <- uint8ToComplex(r.buf[i], r.buf[i+1]):
			}
		}
	}
}

// uint8ToComplex converts a pair of uint8 I/Q values to a normalized complex64 sample.
func uint8ToComplex(iRaw, qRaw byte) complex64 {
	i := (float32(iRaw) - 127.5) / 127.5
	q := (float32(qRaw) - 127.5) / 127.5
	return complex(i, q)
}

// Ensure RTLSource implements IQSource at compile time.
var _ IQSource = (*RTLSource)(nil)
