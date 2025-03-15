package filter

import "github.com/dayvillefire/pocsag-monitor/obj"

func init() {
	RegisterFilter("dummy", func() Filter { return &DummyFilter{} })
}

type DummyFilter struct {
}

func (d *DummyFilter) Init() error {
	return nil
}

// Configure passes the parameters specified to the filter
func (d *DummyFilter) Configure(map[string]any) error {
	return nil
}

// Filter looks at obj.AlphaMessage and processes it
func (d *DummyFilter) Filter(o obj.AlphaMessage) (obj.AlphaMessage, error) {
	return o, nil
}
