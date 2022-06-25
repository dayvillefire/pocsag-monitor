package obj

import "time"

type AlphaMessage struct {
	Timestamp time.Time
	CapCode   string
	Message   string
	Valid     bool
}
