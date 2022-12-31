package output

import (
	"fmt"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

func init() {
	RegisterOutput("stdout", func() Output { return &StdoutOutput{} })
}

type StdoutOutput struct {
}

func (s *StdoutOutput) Init(v string) error {
	return nil
}

func (s *StdoutOutput) SendMessage(a obj.AlphaMessage, channel, msg string) (string, error) {
	fmt.Printf("MESSAGE : %s\n", msg)
	return "", nil
}
