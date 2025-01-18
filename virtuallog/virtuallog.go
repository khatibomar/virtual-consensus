package virtuallog

import "github.com/khatibomar/virtual-consensus/loglet"

type VirtualLoger[T any] interface {
	loglet.Loglet[T]
	Reconfigure() error
}

type VirtualLog[T any] struct {
}

func (v *VirtualLog[T]) Reconfigure() error {
	return nil
}
