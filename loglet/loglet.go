package loglet

import "errors"

var ErrOutOfBounds = errors.New("out of bounds")
var ErrSealed = errors.New("loglet is sealed")

type Loglet[T any] interface {
	Append(Entry T) (int64, error)
	CheckTail() int64
	ReadNext(start, end int64) ([]T, error)
	Seal()
	String() string
}
