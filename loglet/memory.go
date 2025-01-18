package loglet

import "sync/atomic"

type MemoryLoglet[T any] struct {
	entries []T
	sealed  atomic.Bool
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
	if start > end || start < 0 || start >= int64(len(m.entries)) || end < 0 {
		return nil, ErrOutOfBounds
	}

	if end >= int64(len(m.entries)) {
		end = int64(len(m.entries)) - 1
	}

	return m.entries[start : end+1], nil
}

func (m *MemoryLoglet[T]) Seal() {
	m.sealed.Store(true)
}
