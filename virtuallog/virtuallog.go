package virtuallog

import "github.com/khatibomar/virtual-consensus/loglet"

type VirtualLoger[T any] interface {
	loglet.Loglet[T]
	Reconfigure() error
}

type VirtualLog[T any] struct {
	m *MetaStore[T]
}

func NewVirtualLog[T any]() *VirtualLog[T] {
	m := NewMetaStore[T]()
	return &VirtualLog[T]{m: m}
}

func (v *VirtualLog[T]) Append(value T) (int64, error) {
	if v.m.reconfiguring.Load() {
		return 0, ErrReconfiguring
	}
	latest := v.m.latestInChain()
	pos, err := v.m.loglets[latest.LogletID].Append(value)
	if err != nil {
		return 0, err
	}
	return pos + latest.Range.Start, nil
}

func (v *VirtualLog[T]) CheckTail() int64 {
	latest := v.m.latestInChain()
	return latest.Range.Start + v.m.loglets[latest.LogletID].CheckTail()
}

func (v *VirtualLog[T]) Reconfigure() error {
	m := v.m
	return m.Reconfigure(m.version.Load())
}

func (v *VirtualLog[T]) ReadNext(start, end int64) ([]T, error) {
	if start > end || start < 0 || start >= int64(v.CheckTail()) || end < 0 {
		return nil, loglet.ErrOutOfBounds
	}

	if end > int64(v.CheckTail()) {
		end = int64(v.CheckTail())
	}

	result := make([]T, 0, end-start+1)

	firstChain := &v.m.chain
	for start > firstChain.Range.End {
		firstChain = firstChain.Next
	}
	lastChain := firstChain
	for end > firstChain.Range.End {
		lastChain = lastChain.Next
	}

	for {
		log := v.m.loglets[firstChain.LogletID]
		entries, err := log.ReadNext(start-firstChain.Range.Start, log.CheckTail())
		if err != nil {
			return nil, err
		}
		result = append(result, entries...)
		start = 0
		if firstChain == lastChain {
			break
		}
		firstChain = firstChain.Next
	}

	return result, nil
}

func (v *VirtualLog[T]) Seal() {
	m := v.m.loglets[v.m.latestInChain().LogletID]
	m.Seal()
}
