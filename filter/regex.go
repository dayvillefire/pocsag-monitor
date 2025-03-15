package filter

import (
	"fmt"
	"log"
	"regexp"
	"slices"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

func init() {
	RegisterFilter("dummy", func() Filter { return &DummyFilter{} })
}

type RegexFilter struct {
	Regexes []string

	comp []*regexp.Regexp
}

func (d *RegexFilter) Init() error {
	d.comp = make([]*regexp.Regexp, 0)

	return nil
}

// Configure passes the parameters specified to the filter
func (d *RegexFilter) Configure(cfg map[string]any) error {
	_, ok := cfg["regexps"]
	if !ok {
		return fmt.Errorf("no configuration found")
	}
	regexps, ok := cfg["regexps"].([]string)
	if !ok {
		return fmt.Errorf("no configuration found")
	}

	d.Regexes = regexps

	for _, r := range d.Regexes {
		re, err := regexp.Compile(r)
		if err != nil {
			log.Printf("ERR: RexexpFilter: %s", err.Error())
			return err
		}
		d.comp = append(d.comp, re)
	}
	return nil
}

// Filter looks at obj.AlphaMessage and processes it
func (d *RegexFilter) Filter(o obj.AlphaMessage) (obj.AlphaMessage, error) {
	var err error
	x := o.Message
	for _, re := range d.comp {
		x, err = d.process(re, x)
		if err != nil {
			o.Message = x
			return o, err
		}
	}
	o.Message = x
	return o, nil
}

func (d *RegexFilter) process(re *regexp.Regexp, s string) (string, error) {
	// If we don't patch, return
	if !re.MatchString(s) {
		return s, nil
	}

	sm := re.FindAllStringSubmatch(s, -1)
	log.Printf("sm = %#v", sm)
	smi := re.FindAllStringSubmatchIndex(s, -1)
	log.Printf("smi = %#v", smi)
	if len(smi) < 1 {
		return s, nil
	}
	log.Printf("loop")
	for i, x := range slices.Backward(smi) {
		log.Printf("in loop : %d / %v", i, x)
		// Don't process entire string
		if i == 0 {
			continue
		}
		//s = strings.Replace(s, sm[i], "", x)
	}
	return s, nil
}
