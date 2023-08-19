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

type fiber[In any, Out any] struct {
	status uint32
	fn     func(In, SuspendFunc[In, Out]) Out
	ret    Out
	in     chan In
	out    chan Out
}

type Fiber[In any, Out any] interface {
	Start(In) (Out, error)
	Resume(In) (Out, error)
	GetReturn() (Out, error)

	IsStarted() bool
	IsRunning() bool
	IsSuspended() bool
	IsTerminated() bool
}

type SuspendFunc[In any, Out any] func(Out) In

func New[In any, Out any](fn func(In, SuspendFunc[In, Out]) Out) Fiber[In, Out] {
	return &fiber[In, Out]{
		fn:  fn,
		in:  make(chan In),
		out: make(chan Out),
	}
}

func (f *fiber[In, Out]) Start(in In) (Out, error) {
	if atomic.SwapUint32(&f.status, Started|Running) != 0 {
		var zero Out
		return zero, fmt.Errorf("%w: fiber already started", ErrFiber)
	}
	go func() {
		f.ret = f.fn(in, f.suspend)
		close(f.in)
		close(f.out)
		f.status = 0 // Terminated
	}()
	return <-f.out, nil
}

func (f *fiber[In, Out]) Resume(in In) (Out, error) {
	if atomic.SwapUint32(&f.status, f.status|Running) != Started {
		var zero Out
		return zero, fmt.Errorf("%w: fiber not suspend", ErrFiber)
	}
	f.in <- in
	return <-f.out, nil
}

func (f *fiber[In, Out]) suspend(v Out) In {
	f.status &^= Running
	f.out <- v
	return <-f.in
}

func (f *fiber[In, Out]) GetReturn() (Out, error) {
	if f.IsStarted() {
		var zero Out
		return zero, fmt.Errorf("%w: fiber not return", ErrFiber)
	}
	return f.ret, nil
}

func (f *fiber[In, Out]) IsStarted() bool {
	return f.status&Started > 0
}

func (f *fiber[In, Out]) IsRunning() bool {
	return f.status&Running > 0
}

func (f *fiber[In, Out]) IsSuspended() bool {
	return !f.IsRunning()
}

func (f *fiber[In, Out]) IsTerminated() bool {
	return !f.IsStarted()
}
