package fiber

import (
	"sync/atomic"
)

const (
	Started = 1 << iota
	Running
)

type fiber struct {
	status uint32
	fn     func(FiberSuspend)
	in     chan interface{}
	out    chan interface{}
}

type Fiber interface {
	Start() interface{}
	Resume(v interface{}) interface{}

	IsStarted() bool
	IsRunning() bool
	IsSuspended() bool
	IsTerminated() bool
}

type FiberSuspend interface {
	Suspend(v interface{}) interface{}
}

func New(fn func(FiberSuspend)) Fiber {
	return &fiber{
		fn:  fn,
		in:  make(chan interface{}),
		out: make(chan interface{}),
	}
}

func (f *fiber) Start() (v interface{}) {
	if atomic.SwapUint32(&f.status, Started|Running) != 0 {
		return nil
	}
	go func() {
		f.fn(f)
		close(f.in)
		close(f.out)
		f.status = 0 // Terminated
	}()
	return <-f.out
}

func (f *fiber) Resume(v interface{}) interface{} {
	if atomic.SwapUint32(&f.status, f.status|Running) != Started {
		return nil
	}
	f.in <- v
	return <-f.out
}

func (f *fiber) Suspend(v interface{}) interface{} {
	f.status &^= Running
	f.out <- v
	return <-f.in
}

func (f *fiber) IsStarted() bool {
	return f.status&Started > 0
}

func (f *fiber) IsRunning() bool {
	return f.status&Running > 0
}

func (f *fiber) IsSuspended() bool {
	return !f.IsRunning()
}

func (f *fiber) IsTerminated() bool {
	return !f.IsStarted()
}
