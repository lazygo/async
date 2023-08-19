package fiber

import (
	"fmt"
	"testing"
)

func TestFiber(t *testing.T) {
	fi := New(func(in string, suspend SuspendFunc[string, string]) string {
		if in != "fiber start" {
			fmt.Println(in)
			t.Error("Start error")
		}
		v := suspend("first suspend")
		if v != "second resume" {
			fmt.Println(v)
			t.Error("Resume error")
		}
		return "return val"
	})

	first, err := fi.Start("fiber start")
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
