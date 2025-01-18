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
		if err != ErrOutOfBounds {
			t.Errorf("Expected ErrOutOfBounds, got %v", err)
		}

		_, err = vlog.ReadNext(0, -1)
		if err != ErrOutOfBounds {
			t.Errorf("Expected ErrOutOfBounds, got %v", err)
		}

		_, err = vlog.ReadNext(1, 0)
		if err != ErrOutOfBounds {
			t.Errorf("Expected ErrOutOfBounds, got %v", err)
		}
	})

	t.Run("read from many loglets", func(t *testing.T) {
		vlog := NewVirtualLog[string]()
		_, _ = vlog.Append("first")
		_, _ = vlog.Append("second")
		err := vlog.Reconfigure()
		if err != nil {
			t.Errorf("Reconfigure failed: %v", err)
		}
		pos, err := vlog.Append("third")
		if err != nil {
			t.Errorf("Append failed: %v", err)
		}
		if pos != 2 {
			t.Errorf("Expected position 2, got %d", pos)
		}
		pos, err = vlog.Append("fourth")
		if err != nil {
			t.Errorf("Append failed: %v", err)
		}
		if pos != 3 {
			t.Errorf("Expected position 3, got %d", pos)
		}

		entries, err := vlog.ReadNext(0, 4)
		if err != nil {
			t.Errorf("ReadNext failed: %v", err)
		}
		if len(entries) != 4 {
			t.Errorf("Expected 4 entries, got %d", len(entries))
		}
		if entries[0] != "first" || entries[1] != "second" || entries[2] != "third" || entries[3] != "fourth" {
			t.Errorf("Unexpected entries: %v", entries)
		}

		entries, err = vlog.ReadNext(1, 3)
		if err != nil {
			t.Errorf("ReadNext failed: %v", err)
		}
		if len(entries) != 3 {
			t.Errorf("Expected 3 entries, got %d", len(entries))
		}
		if entries[0] != "second" || entries[1] != "third" || entries[2] != "fourth" {
			t.Errorf("Unexpected entries: %v", entries)
		}

		entries, err = vlog.ReadNext(0, 2)
		if err != nil {
			t.Errorf("ReadNext failed: %v", err)
		}
		if len(entries) != 3 {
			t.Errorf("Expected 3 entries, got %d", len(entries))
		}
		if entries[0] != "first" || entries[1] != "second" || entries[2] != "third" {
			t.Errorf("Unexpected entries: %v", entries)
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
