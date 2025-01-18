package virtuallog

import (
	"testing"

	"github.com/khatibomar/virtual-consensus/loglet"
)

func TestVirtualLog(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		vlog := NewVirtualLog[string]()

		pos1, err := vlog.Append("first")
		if err != nil {
			t.Errorf("Append failed: %v", err)
		}
		if pos1 != 0 {
			t.Errorf("Expected position 0, got %d", pos1)
		}

		pos2, err := vlog.Append("second")
		if err != nil {
			t.Errorf("Append failed: %v", err)
		}
		if pos2 != 1 {
			t.Errorf("Expected position 1, got %d", pos2)
		}

		tail := vlog.CheckTail()
		if tail != 1 {
			t.Errorf("Expected tail 1, got %d", tail)
		}

		entries, err := vlog.ReadNext(0, 1)
		if err != nil {
			t.Errorf("ReadNext failed: %v", err)
		}
		if len(entries) != 2 {
			t.Errorf("Expected 2 entries, got %d", len(entries))
		}
		if entries[0] != "first" || entries[1] != "second" {
			t.Errorf("Unexpected entries: %v", entries)
		}
	})

	t.Run("out of bounds read", func(t *testing.T) {
		vlog := NewVirtualLog[string]()
		_, err := vlog.ReadNext(-1, 0)
		if err != loglet.ErrOutOfBounds {
			t.Errorf("Expected ErrOutOfBounds, got %v", err)
		}

		_, err = vlog.ReadNext(0, -1)
		if err != loglet.ErrOutOfBounds {
			t.Errorf("Expected ErrOutOfBounds, got %v", err)
		}

		_, err = vlog.ReadNext(1, 0)
		if err != loglet.ErrOutOfBounds {
			t.Errorf("Expected ErrOutOfBounds, got %v", err)
		}
	})

	t.Run("seal", func(t *testing.T) {
		vlog := NewVirtualLog[string]()
		p, err := vlog.Append("test")
		if err != nil {
			t.Errorf("Append failed: %v", err)
		}
		if p != 0 {
			t.Errorf("Expected position 0, got %d", p)
		}
		vlog.Seal()
		_, err = vlog.Append("test")
		if err != loglet.ErrSealed {
			t.Errorf("Expected ErrSealed, got %v", err)
		}
	})
}
