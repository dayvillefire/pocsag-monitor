package filter

import (
	"strings"

	"github.com/dayvillefire/pocsag-monitor/obj"
	"github.com/jbuchbinder/shims"
)

func init() {
	RegisterFilter("splitmatch", func() Filter { return &SplitMatchFilter{} })
}

type SplitMatchFilter struct {
	ifmatch string
	split   string
	fields  []int
}

func (d *SplitMatchFilter) Init() error {
	return nil
}

// Configure passes the parameters specified to the filter
func (d *SplitMatchFilter) Configure(cfg map[string]any) error {
	ifmatch, err := ConfigValue[string](cfg, "if-match")
	if err != nil {
		ifmatch = ""
	}

	d.ifmatch = ifmatch

	split, err := ConfigValue[string](cfg, "split")
	if err != nil {
		return err
	}

	d.split = split

	fields, err := ConfigValue[[]int](cfg, "remove-fields")

	if err != nil {
		return err
	}

	d.fields = fields

	return nil
}

// Filter looks at obj.AlphaMessage and processes it
func (d *SplitMatchFilter) Filter(o obj.AlphaMessage) (obj.AlphaMessage, error) {
	msg := o.Message

	// If there's a matching string, check for it
	if d.ifmatch != "" && !strings.Contains(msg, d.ifmatch) {
		return o, nil
	}

	// If there is no split in the string, return as is
	if !strings.Contains(msg, d.split) {
		return o, nil
	}

	parts := strings.Split(msg, d.split)
	f := []string{}
	for pos, part := range parts {
		if shims.InArray(pos, d.fields) {
			continue
		}
		f = append(f, part)
	}

	o.Message = strings.Join(f, d.split)

	return o, nil
}
