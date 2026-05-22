package sdr

import (
	"io"
	"os"
)

// FileIQSource replays IQ samples from a file for testing.
// File format: raw interleaved uint8 I/Q pairs [I0, Q0, I1, Q1, ...].
type FileIQSource struct {
	path   string
	ch     chan IQSample
	done   chan struct{}
	reader *os.File
}

func NewFileIQSource(path string) *FileIQSource {
	return &FileIQSource{path: path}
}

func (f *FileIQSource) Start(_ int, _ int, _ int, _ int) (<-chan IQSample, error) {
	var err error
	f.reader, err = os.Open(f.path)
	if err != nil {
		return nil, err
	}
	f.ch = make(chan IQSample, 4096)
	f.done = make(chan struct{})
	go f.stream()
	return f.ch, nil
}

func (f *FileIQSource) stream() {
	defer close(f.ch)
	defer f.reader.Close()
	buf := make([]byte, 2) // one I/Q pair
	for {
		select {
		case <-f.done:
			return
		default:
		}
		_, err := io.ReadFull(f.reader, buf)
		if err != nil {
			return // EOF or error ends stream
		}
		// uint8 values centered at 127.5 -> normalized to [-1, 1]
		i := (float32(buf[0]) - 127.5) / 127.5
		q := (float32(buf[1]) - 127.5) / 127.5
		f.ch <- complex(i, q)
	}
}

func (f *FileIQSource) Stop() error {
	close(f.done)
	return nil
}

// Ensure FileIQSource implements IQSource.
var _ IQSource = (*FileIQSource)(nil)
