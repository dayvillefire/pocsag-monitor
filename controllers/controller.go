package controllers

import (
	"github.com/jbuchbinder/shims/factory"
)

type Controller interface {
	Run() error
	Stop() error
}

type ControllerDummy struct {
}

func (d ControllerDummy) Run() error  { return nil }
func (d ControllerDummy) Stop() error { return nil }

var (
	ControllerFactory factory.Factory[Controller]
)

func init() {
	ControllerFactory = factory.New[Controller](ControllerDummy{})
}
