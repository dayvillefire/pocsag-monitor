package output

import (
	"github.com/dayvillefire/pocsag-monitor/obj"
)

func init() {
	RegisterOutput("dummy", func() Output { return &DummyOutput{} })
}

type DummyOutput struct {
}

func (s *DummyOutput) Init(v string) error {
	return nil
}

func (s *DummyOutput) SendMessage(a obj.AlphaMessage, channel, msg string) (string, error) {
	return "", nil
}
