package tq

import (
	"fmt"
	"testing"
	"time"
)

func TestTq(t *testing.T) {
	fn := func(k int, v *string) {
		fmt.Printf("fn %d: %s => %s\n", k, *v, time.Now().Format(time.RFC1123))
	}
	tq := NewTimerQueue[int, string](fn, func() time.Duration { return time.Duration(3 * time.Second) })
	tq.CheckConsistency()
	tq.Push(1, "aaa")
	tq.CheckConsistency()
	time.Sleep(time.Second)
	tq.Push(2, "bbb")
	tq.CheckConsistency()
	time.Sleep(time.Second)
	tq.Push(3, "ccc")
	tq.CheckConsistency()
	if tq.Remove(2) {
		fmt.Printf("item 2 removed\n")
	}
	tq.CheckConsistency()
	time.Sleep(time.Second)
	tq.Refresh(3)
	tq.CheckConsistency()
	time.Sleep(1 * time.Second)
	tq.Refresh(3)
	tq.CheckConsistency()
	time.Sleep(1 * time.Second)
	tq.Refresh(3)
	tq.CheckConsistency()
	time.Sleep(1 * time.Second)
	tq.Refresh(3)
	tq.CheckConsistency()
	time.Sleep(1 * time.Second)
	tq.Refresh(3)
	tq.CheckConsistency()
	time.Sleep(1 * time.Second)
	tq.Refresh(3)
	time.Sleep(1 * time.Second)
	tq.CheckConsistency()
	tq.Refresh(3)

	time.Sleep(5 * time.Second)
	//for {
	//	if tq.IsEmpty() {
	//		break
	//	}
	//	time.Sleep(100 * time.Millisecond)
	//}
}

func TestTq2(t *testing.T) {
	assert := func(cond bool, msg string) {
		if !cond {
			t.Fail()
		}
	}
	fn := func(k int, v *string) {
		fmt.Printf("fn %d: %s => %s\n", k, *v, time.Now().Format(time.RFC1123))
	}
	tq := NewTimerQueue[int, string](fn, func() time.Duration { return time.Duration(3 * time.Second) })
	tq.Push(1, "a")
	tq.Remove(1)
	tq.CheckConsistency()
	assert(tq.Get(1) == nil, "1")
	tq.Push(1, "a+")
	tq.CheckConsistency()
	assert(*tq.Get(1) == "a+", "2")
	tq.Push(2, "b")
	tq.CheckConsistency()
	tq.Push(3, "c")
	tq.CheckConsistency()
	tq.Push(4, "d")
	tq.CheckConsistency()
	tq.Push(5, "e")
	tq.CheckConsistency()
	tq.Remove(1)
	tq.CheckConsistency()
	tq.Remove(4)
	tq.CheckConsistency()
	tq.Remove(5)
	tq.CheckConsistency()
	tq.Remove(3)
	tq.CheckConsistency()
	time.Sleep(5 * time.Second)
}
