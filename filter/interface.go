package filter

import (
	"errors"
	"fmt"
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
	Init() error

	// Configure passes the parameters specified to the filter
	Configure(map[string]any) error

	// Filter looks at obj.AlphaMessage and processes it
	Filter(obj.AlphaMessage) (obj.AlphaMessage, error)
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
func RegisterFilter(name string, o func() Filter) {
	filterMapLock.Lock()
	defer filterMapLock.Unlock()
	log.Printf("RegisterFilter: %s", name)
	filterMap[name] = o
}

func ConfigValue[T any](cfg map[string]any, key string) (T, error) {
	_, ok := cfg[key]
	if !ok {
		return *(new(T)), fmt.Errorf("key not present")
	}
	v, ok := cfg[key].(T)
	if !ok {
		return *(new(T)), fmt.Errorf("wrong type")
	}
	return v, nil
}
