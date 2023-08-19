package fiber

import (
	"fmt"
	"testing"
)

func TestFiber(t *testing.T) {
	fi := New(func(suspend SuspendFunc[string]) string {
		fmt.Println("started")
		v := suspend("first suspend")
		if v != "second resume" {
			fmt.Println(v)
			t.Error("Resume error")
		}
		return "return val"
	})

	first, err := fi.Start()
	if err != nil {
		t.Error(err)
	}
	if first != "first suspend" {
		fmt.Println(first)
		t.Error("Suspend error")
	}
	last, err := fi.Resume("second resume")
	if err != nil {
		t.Error(err)
	}
	if last != "" {
		fmt.Println(last)
		t.Error("Terminated error")
	}

	ret, err := fi.GetReturn()
	if err != nil {
		t.Error(err)
	}
	if ret != "return val" {
		fmt.Println(ret)
		t.Error("Terminated error")
	}

}

func TestFiber1(t *testing.T) {

	fib := New(func(suspend SuspendFunc[string]) string {
		value := suspend("fiber")
		fmt.Println("Value used to resume fiber:", value)
		return "return val"
	})

	value, _ := fib.Start()
	fmt.Println("Value from fiber suspending:", value)
	fib.Resume("test")
	ret, _ := fib.GetReturn()
	fmt.Println("Value from fiber return:", ret)
}
