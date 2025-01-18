package loglet

import "testing"

func TestMemoryLoglet(t *testing.T) {
	t.Run("creation", func(t *testing.T) {
		m := NewMemoryLoglet[int]()
		if m == nil {
			t.Fatal("MemoryLoglet is nil")
		}
	})

	t.Run("append and get", func(t *testing.T) {
		m := NewMemoryLoglet[int]()
		pos, err := m.Append(1)
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
		if pos != 0 {
			t.Errorf("expected 0, got %d", pos)
		}
		_, _ = m.Append(2)
		_, _ = m.Append(3)

		if m.CheckTail() != 2 {
			t.Errorf("expected tail to be 2, got %d", m.CheckTail())
		}

		logs := m.entries
		if len(logs) != 3 {
			t.Errorf("expected 3 logs, got %d", len(logs))
		}
		if logs[0] != 1 || logs[1] != 2 || logs[2] != 3 {
			t.Error("logs not stored in correct order")
		}
	})

	t.Run("seal", func(t *testing.T) {
		m := NewMemoryLoglet[int]()
		_, _ = m.Append(1)
		_, _ = m.Append(2)
		m.Seal()

		if !m.sealed.Load() {
			t.Error("expected loglet to be sealed")
		}

		_, err := m.Append(3)
		if err != ErrSealed {
			t.Errorf("expected ErrSealed, got %v", err)
		}

		if len(m.entries) != 2 {
			t.Errorf("expected 2 logs, got %d", len(m.entries))
		}
	})

	t.Run("concurrent operations", func(t *testing.T) {
		m := NewMemoryLoglet[int]()
		done := make(chan bool)

		go func() {
			for i := 0; i < 100; i++ {
				_, _ = m.Append(i)
			}
			done <- true
		}()

		go func() {
			for i := 100; i < 200; i++ {
				_, _ = m.Append(i)
			}
			done <- true
		}()

		<-done
		<-done

		logs := m.entries
		if len(logs) != 200 {
			t.Errorf("expected 200 logs, got %d", len(logs))
		}
	})
}
