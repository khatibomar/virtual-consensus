package virtuallog

import (
	"fmt"
	"strings"
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

func (c *Chain) String() string {
	sb := strings.Builder{}
	sb.WriteString("[------- ")
	sb.WriteString(c.LogletID)
	sb.WriteString(" -------]\n")
	sb.WriteString("Range: ")
	sb.WriteString("[")
	start := c.Range.Start
	sb.WriteString(fmt.Sprintf("%v", start))
	sb.WriteString(", ")
	end := c.Range.End
	if end == Infinity {
		sb.WriteString("∞")
		sb.WriteString(")")
	} else {
		sb.WriteString(fmt.Sprintf("%v", end))
		sb.WriteString("]")
	}
	sb.WriteString("\n")
	return sb.String()
}

type MetaStore[T any] struct {
	loglets map[string]loglet.Loglet[T]
	chain   Chain

	version       *atomic.Int32
	keyGenerator  func() string
	reconfiguring *atomic.Bool
}

func (m *MetaStore[T]) String() string {
	chain := &m.chain
	sb := strings.Builder{}
	for chain != nil {
		sb.WriteString(chain.String())
		loglet := m.loglets[chain.LogletID]
		sb.WriteString(loglet.String())
		sb.WriteString("\n")
		chain = chain.Next
	}
	return sb.String()
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
