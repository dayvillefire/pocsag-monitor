package output

import (
	"errors"
	"sync"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

var (
	outputMap     = map[string]func() Output{}
	outputMapLock = new(sync.Mutex)
)

type Output interface {
	// Init initializes the plugin
	Init(string) error
	// SendMessage specifies an obj.AlphaMessage object, a plugin option or
	// channel, and the message text to be sent
	SendMessage(obj.AlphaMessage, string, string) (string, error)
}

// InstantiateOutput instantiates an Output by name
func InstantiateOutput(name string) (o Output, err error) {
	var f func() Output
	var found bool
	if f, found = outputMap[name]; !found {
		err = errors.New("unable to resolve output " + name)
		return
	}
	o = f()
	err = nil
	return
}

// RegisterOutput adds a new Output instance to the registry
func RegisterOutput(name string, o func() Output) {
	outputMapLock.Lock()
	defer outputMapLock.Unlock()
	outputMap[name] = o
}
