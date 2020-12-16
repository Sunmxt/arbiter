package arbiter

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArbiter(t *testing.T) {

	spawnRoutines := func(a *Arbiter, n uint32, before func()) {
		for idx := uint32(0); idx < n; idx++ {
			a.Go(func() {
				<-a.Exit()
				if before != nil {
					before()
				}
			})
		}
	}

	t.Run("new_normal", func(t *testing.T) {
		t.Parallel()
		a := New()
		assert.NotNil(t, a)
		assert.True(t, a.ShouldRun())
	})

	t.Run("new_with_parent", func(t *testing.T) {
		t.Parallel()
		p := New()
		assert.NotNil(t, p)
		assert.True(t, p.ShouldRun())

		a := NewWithParent(p)
		assert.NotNil(t, a)
		assert.True(t, a.ShouldRun())
	})

	t.Run("new_with_shut_parent", func(t *testing.T) {
		t.Parallel()
		p := New()
		assert.NotNil(t, p)
		assert.True(t, p.ShouldRun())
		p.Shutdown()
		assert.False(t, p.ShouldRun())

		a := NewWithParent(p)
		assert.NotNil(t, a)
		assert.False(t, a.ShouldRun())
	})

	t.Run("test_exit_order_no_parent", func(t *testing.T) {
		t.Parallel()
		p := New()

		exitCounter := uint32(10)

		spawnRoutines(p, exitCounter, func() {
			atomic.AddUint32(&exitCounter, 0xFFFFFFFF)
		})

		p.Shutdown()
		p.Join()
		assert.Equal(t, uint32(0xFFFFFFFF), atomic.AddUint32(&exitCounter, 0xFFFFFFFF))
	})

	t.Run("test_tree", func(t *testing.T) {
		t.Parallel()
		p := New()
		l, r := NewWithParent(p), NewWithParent(p)

		spawnRoutines(p, 10, nil)
		spawnRoutines(l, 9, nil)
		spawnRoutines(r, 11, nil)

		assert.Equal(t, int32(9), l.NumGoroutine())
		assert.Equal(t, int32(11), r.NumGoroutine())
		assert.Equal(t, int32(30), p.NumGoroutine())

		assert.True(t, l.ShouldRun())
		assert.True(t, r.ShouldRun())
		p.Shutdown()
		assert.False(t, l.ShouldRun())
		assert.False(t, r.ShouldRun())

		p.Join()
		l.Join()
		r.Join()
	})

	t.Run("test_tick_go_surge", func(t *testing.T) {
		t.Parallel()
		p := New()

		concurrency := int32(0)
		p.TickGo(func(cancel func(), deadline time.Time) {
			assert.LessOrEqual(t, atomic.AddInt32(&concurrency, 1), int32(3))
			time.Sleep(time.Millisecond * 600)
			atomic.AddInt32(&concurrency, -1)
		}, time.Millisecond*100, 3)

		time.Sleep(time.Second * 3)
		p.Shutdown()
		p.Join()
	})

	t.Run("test_tick_go_period", func(t *testing.T) {
		t.Parallel()
		p := New()

		var last time.Time

		period, first := time.Millisecond*100, true

		p.TickGo(func(cancel func(), deadline time.Time) {
			if !first {
				assert.WithinDuration(t, last.Add(period), deadline, time.Millisecond*5)
			}
			last, first = deadline, false
		}, time.Millisecond*100, 1)

		time.Sleep(time.Second)
		p.Shutdown()
		p.Join()
	})

	t.Run("hooks", func(t *testing.T) {
		p := New()

		prestopped, stopped := false, false
		prestop := func() { prestopped = true }
		stop := func() { stopped = true }

		p.HookPreStop(prestop).HookStopped(stop)
		spawnRoutines(p, 10, nil)

		p.Shutdown()
		p.Join()

		assert.True(t, prestopped)
		assert.True(t, stopped)
	})
}
