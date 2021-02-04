package state

import "testing"

func TestState_Format(t *testing.T) {
	states := []State{
		{LongTermRule: 1},
		{LongTermRule: -1},
		{LongTermRule: 0},
	}
	t.Run("file format", func(t *testing.T) {
		for _, s := range states {
			t.Logf(s.Format(false))
		}
	})
	t.Run("file format", func(t *testing.T) {
		for _, s := range states {
			t.Logf(s.Format(true))
		}
	})
}
