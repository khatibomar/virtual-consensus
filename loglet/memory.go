package loglet

import (
	"fmt"
	"strings"
	"sync/atomic"
)

type MemoryLoglet[T any] struct {
	entries []T
	sealed  atomic.Bool
}

func (m *MemoryLoglet[T]) String() string {
	if len(m.entries) == 0 {
		return "No entries\n"
	}
	sb := strings.Builder{}
	sb.WriteString("Entries: \n")
	for _, entry := range m.entries {
		sb.WriteString("\t")
		sb.WriteString(fmt.Sprintf("%v\n", entry))
	}
	return sb.String()
}

func NewMemoryLoglet[T any]() *MemoryLoglet[T] {
	return &MemoryLoglet[T]{
		entries: make([]T, 0),
		sealed:  atomic.Bool{},
	}
}

func (m *MemoryLoglet[T]) Append(entry T) (int64, error) {
	if m.sealed.Load() {
		return 0, ErrSealed
	}
	m.entries = append(m.entries, entry)
	return int64(len(m.entries) - 1), nil
}

func (m *MemoryLoglet[T]) CheckTail() int64 {
	return int64(len(m.entries) - 1)
}

func (m *MemoryLoglet[T]) ReadNext(start, end int64) ([]T, error) {
	if start > end || start < 0 || start >= m.CheckTail() || end < 0 {
		return nil, ErrOutOfBounds
	}

	if end >= m.CheckTail() {
		end = m.CheckTail()
	}

	return m.entries[start : end+1], nil
}

func (m *MemoryLoglet[T]) Seal() {
	m.sealed.Store(true)
}
