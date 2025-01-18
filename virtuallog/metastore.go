package virtuallog

import (
	"fmt"
	"sync/atomic"

	"github.com/khatibomar/virtual-consensus/loglet"
)

var ErrVersionMismatch = fmt.Errorf("version mismatch")
var ErrReconfiguring = fmt.Errorf("reconfiguring")
var Infinity = int64(^uint64(0) >> 1)

type Range struct {
	Start int64
	End   int64
}

type Chain struct {
	Range    Range
	LogletID string
	Next     *Chain
}

type MetaStore[T any] struct {
	loglets map[string]loglet.Loglet[T]
	chain   Chain

	version       *atomic.Int32
	keyGenerator  func() string
	reconfiguring *atomic.Bool
}

func (m *MetaStore[T]) latestInChain() *Chain {
	chain := &m.chain
	for chain.Next != nil {
		chain = chain.Next
	}
	return chain
}

func NewMetaStore[T any]() *MetaStore[T] {
	m := &MetaStore[T]{}
	m.version = &atomic.Int32{}
	m.reconfiguring = &atomic.Bool{}

	keyGenerator := keyGen()
	m.keyGenerator = keyGenerator

	k := keyGenerator()
	m.loglets = make(map[string]loglet.Loglet[T])
	loglet := loglet.NewMemoryLoglet[T]()
	m.loglets[k] = loglet
	chain := Chain{Range: Range{Start: 0, End: Infinity}, LogletID: k, Next: nil}
	m.chain = chain

	return m
}

func (m *MetaStore[T]) Reconfigure(version int32) error {
	if version != m.version.Load() {
		return ErrVersionMismatch
	}

	if m.reconfiguring.Load() {
		return ErrReconfiguring
	}
	m.reconfiguring.Store(true)

	latest := m.latestInChain()
	latestLoglet := m.loglets[latest.LogletID]
	latestLoglet.Seal()
	latest.Range.End = latest.Range.Start + latestLoglet.CheckTail()

	newLoglet := loglet.NewMemoryLoglet[T]()
	k := m.keyGenerator()
	m.loglets[k] = newLoglet
	newChain := Chain{Range: Range{Start: latest.Range.End + 1, End: Infinity}, LogletID: k, Next: nil}
	latest.Next = &newChain

	m.version.Add(1)
	m.reconfiguring.Store(false)

	return nil
}

func keyGen() func() string {
	i := 0
	return func() string {
		s := fmt.Sprintf("memory%d", i)
		i++
		return s
	}
}
