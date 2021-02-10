package mod

import (
	"fmt"
)

type Stat int

var ()

type Module interface {
	Init(cfg string) bool

	Name() string

	Process() Stat
}

func Register(name string, mod Module) error {
}
