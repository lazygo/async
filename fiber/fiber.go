package fiber

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const (
	Started = 1 << iota
	Running
)

var (
	ErrFiber = errors.New("fiber error")
)

type fiber[T any] struct {
	status uint32
	fn     func(SuspendFunc[T]) T
	ret    T
	in     chan T
	out    chan T
}

type Fiber[T any] interface {
	Start() (T, error)
	Resume(v T) (T, error)
	GetReturn() (T, error)

	IsStarted() bool
	IsRunning() bool
	IsSuspended() bool
	IsTerminated() bool
}

type SuspendFunc[T any] func(T) T

func New[T any](fn func(SuspendFunc[T]) T) Fiber[T] {
	return &fiber[T]{
		fn:  fn,
		in:  make(chan T),
		out: make(chan T),
	}
}
func (f *fiber[T]) Start() (T, error) {
	if atomic.SwapUint32(&f.status, Started|Running) != 0 {
		var zero T
		return zero, fmt.Errorf("%w: fiber already started", ErrFiber)
	}
	go func() {
		f.ret = f.fn(f.suspend)
		close(f.in)
		close(f.out)
		f.status = 0 // Terminated
	}()
	return <-f.out, nil
}

func (f *fiber[T]) Resume(v T) (T, error) {
	if atomic.SwapUint32(&f.status, f.status|Running) != Started {
		var zero T
		return zero, fmt.Errorf("%w: fiber not suspend", ErrFiber)
	}
	f.in <- v
	return <-f.out, nil
}

func (f *fiber[T]) suspend(v T) T {
	f.status &^= Running
	f.out <- v
	return <-f.in
}

func (f *fiber[T]) GetReturn() (T, error) {
	if f.IsStarted() {
		var zero T
		return zero, fmt.Errorf("%w: fiber not return", ErrFiber)
	}
	return f.ret, nil
}

func (f *fiber[T]) IsStarted() bool {
	return f.status&Started > 0
}

func (f *fiber[T]) IsRunning() bool {
	return f.status&Running > 0
}

func (f *fiber[T]) IsSuspended() bool {
	return !f.IsRunning()
}

func (f *fiber[T]) IsTerminated() bool {
	return !f.IsStarted()
}
