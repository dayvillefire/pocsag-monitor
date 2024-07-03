package output

import (
	"errors"
	"log"
	"sync"

	"github.com/dayvillefire/pocsag-monitor/obj"
)

var (
	filterMap     = map[string]func() Filter{}
	filterMapLock = new(sync.Mutex)
)

type Filter interface {
	// Init initializes the plugin
	Init(string) error
	// SendMessage specifies an obj.AlphaMessage object, a plugin option or
	// channel, and the message text to be sent
	SendMessage(obj.AlphaMessage, string, string) (string, error)
}

// InstantiateFilter instantiates an Filter by name
func InstantiateFilter(name string) (o Filter, err error) {
	var f func() Filter
	var found bool
	if f, found = filterMap[name]; !found {
		err = errors.New("unable to resolve filter " + name)
		return
	}
	o = f()
	err = nil
	return
}

// RegisterFilter adds a new Filter instance to the registry
func RegisterOutput(name string, o func() Filter) {
	filterMapLock.Lock()
	defer filterMapLock.Unlock()
	log.Printf("RegisterFilter: %s", name)
	filterMap[name] = o
}
