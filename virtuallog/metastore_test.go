package virtuallog

import (
	"testing"

	"github.com/khatibomar/virtual-consensus/loglet"
)

func TestMetastoreReconfigure(t *testing.T) {
	t.Run("initial reconfiguration", func(t *testing.T) {
		m := NewMetaStore[int]()
		log := m.loglets[m.latestInChain().LogletID]
		_, _ = log.Append(1)
		_, _ = log.Append(2)

		currVersion := m.version.Load()

		err := m.Reconfigure(currVersion)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		if m.version.Load() != currVersion+1 {
			t.Errorf("Expected version %v, got %v", currVersion+1, m.version.Load())
		}

		err = m.Reconfigure(currVersion)
		if err != ErrVersionMismatch {
			t.Errorf("Expected %v, got %v", ErrVersionMismatch, err)
		}
	})

	t.Run("chain validation", func(t *testing.T) {
		m := NewMetaStore[int]()
		log := m.loglets[m.latestInChain().LogletID]
		_, _ = log.Append(1)
		_, _ = log.Append(2)

		firstChain := m.latestInChain()
		_ = m.Reconfigure(m.version.Load())

		if len(m.loglets) != 2 {
			t.Errorf("Expected 2 loglets, got %v", len(m.loglets))
		}

		validateChain(t, firstChain, 0, 1, true)

		secondChain := firstChain.Next
		validateChain(t, secondChain, firstChain.Range.End+1, Infinity, false)
	})

	t.Run("append after reconfiguration", func(t *testing.T) {
		m := NewMetaStore[int]()
		log := m.loglets[m.latestInChain().LogletID]
		_, _ = log.Append(1)
		_, _ = log.Append(2)

		_ = m.Reconfigure(m.version.Load())

		_, err := log.Append(3)
		if err != loglet.ErrSealed {
			t.Errorf("Expected %v, got %v", loglet.ErrSealed, err)
		}

		log = m.loglets[m.latestInChain().LogletID]
		pos, err := log.Append(3)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
		if pos != 0 {
			t.Errorf("Expected position 0, got %v", pos)
		}
	})

	t.Run("multiple reconfigurations", func(t *testing.T) {
		m := NewMetaStore[int]()
		log := m.loglets[m.latestInChain().LogletID]
		_, _ = log.Append(1)
		_, _ = log.Append(2)

		_ = m.Reconfigure(m.version.Load())
		log = m.loglets[m.latestInChain().LogletID]
		_, _ = log.Append(3)

		secondChain := m.latestInChain()
		err := m.Reconfigure(m.version.Load())
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		thirdChain := m.latestInChain()

		if secondChain.Range.End != 2 {
			t.Errorf("Expected second chain end 2, got %v", secondChain.Range.End)
		}

		if thirdChain.Range.Start != 3 {
			t.Errorf("Expected third chain start 3, got %v", thirdChain.Range.Start)
		}
	})
}

func validateChain(t *testing.T, chain *Chain, expectedStart, expectedEnd int64, hasNext bool) {
	t.Helper()

	if hasNext && chain.Next == nil {
		t.Error("Expected chain to have next, got nil")
	} else if !hasNext && chain.Next != nil {
		t.Error("Expected chain to not have next, got non-nil")
	}

	if chain.Range.Start != expectedStart {
		t.Errorf("Expected start %v, got %v", expectedStart, chain.Range.Start)
	}

	if chain.Range.End != expectedEnd {
		t.Errorf("Expected end %v, got %v", expectedEnd, chain.Range.End)
	}
}
