package tq

import (
	"sync"
	"sync/atomic"
	"time"
)

type tqItem[K, V any] struct {
	next    *tqItem[K, V]
	prev    *tqItem[K, V]
	key     K
	value   V
	expires atomic.Int64
}

type TimerQueue[K comparable, V any] struct {
	front *tqItem[K, V]
	back  *tqItem[K, V]
	fn    func(K, *V)
	dur   func() time.Duration
	m     map[K]*tqItem[K, V]
	sync.Mutex
}

// Create new timer queue with items accessible by key of type K and values of type V.
func NewTimerQueue[K comparable, V any](expire_func func(K, *V), duration func() time.Duration) *TimerQueue[K, V] {
	tq := TimerQueue[K, V]{front: nil, back: nil, fn: expire_func, dur: duration, m: make(map[K]*tqItem[K, V])}
	return &tq
}

// returns true if queue was empty
func (tq *TimerQueue[K, V]) push_back(item *tqItem[K, V]) bool {
	if nil == tq.front {
		tq.front, tq.back = item, item
		return true
	} else {
		item.prev = tq.back // item.next = nil
		tq.back.next = item
		tq.back = item
		return false
	}
}

func (tq *TimerQueue[K, V]) remove(item *tqItem[K, V]) {
	if item == tq.front {
		tq.pop_front()
		return
	}
	if item.prev == nil { // avoid removal twice
		return
	}
	item.prev.next = item.next
	if item.next != nil {
		item.next.prev = item.prev
	} else {
		tq.back = item.prev
	}
}

func (tq *TimerQueue[K, V]) pop_front() *tqItem[K, V] {
	if nil == tq.front {
		return nil
	}
	tmp := tq.front
	if tq.front.next == nil { // single item
		tq.front, tq.back = nil, nil
	} else {
		tq.front.next.prev = nil
		tq.front = tq.front.next
	}
	return tmp
}

func (tq *TimerQueue[K, V]) set_expire(i *tqItem[K, V]) {
	i.expires.Store(time.Now().Add(tq.dur()).UnixNano())
}

// Push new item at the tail of the timer queue. Item will expire at the current time plus duration returned
// by 'duration' argument of NewTimerQueue().
func (tq *TimerQueue[K, V]) Push(k K, v V) {
	i := &tqItem[K, V]{next: nil, key: k, value: v}
	tq.set_expire(i)
	tq.Lock()
	defer tq.Unlock()
	tq.m[k] = i
	if tq.push_back(i) {
		go tq.runTimer()
	}
}

// Refresh expiration time of an item (re-insert it into the back of the queue)
func (tq *TimerQueue[K, V]) Refresh(k K) {
	tq.Lock()
	defer tq.Unlock()
	if item := tq.m[k]; nil != item {
		tq.remove(item)
		tq.push_back(item)
		tq.set_expire(item)
	}
}

// Get item by it's key. Returns nil if item not found
func (tq *TimerQueue[K, V]) Get(k K) *V {
	tq.Lock()
	defer tq.Unlock()
	if item := tq.m[k]; nil != item {
		return &item.value
	}
	return nil
}

// Remove item. Returns true if item existed and was removed
func (tq *TimerQueue[K, V]) Remove(k K) bool {
	tq.Lock()
	defer tq.Unlock()
	if item := tq.m[k]; nil != item {
		delete(tq.m, k)
		tq.remove(item)
		return true
	}
	return false
}

// Returns true if timer queue is empty
func (tq *TimerQueue[K, V]) IsEmpty() bool {
	tq.Lock()
	defer tq.Unlock()
	return tq.front == nil
}

func (tq *TimerQueue[K, V]) runTimer() {
	for {
		tq.Lock()
		item := tq.front
		if nil == item {
			tq.Unlock()
			return // exit loop
		}
		tq.Unlock()
		ext := item.expires.Load()
		t := ext - time.Now().UnixNano()
		if t > 0 {
			time.Sleep(time.Duration(t) * time.Nanosecond)
		}
		tq.Lock()
		if item != tq.front || item.expires.Load() != ext { // refresh happened or removed meanwhile
			tq.Unlock()
			continue // repeat
		}
		delete(tq.m, item.key)
		tq.remove(item)
		tq.Unlock()
		tq.fn(item.key, &item.value)
	}
}

// for testing purposes only
func (tq *TimerQueue[K, V]) CheckConsistency() {
	assert := func(cond bool, msg string) {
		if !cond {
			panic(msg)
		}
	}
	if tq.front == nil {
		assert(tq.back == nil, "1")
		return
	}
	if tq.front.next == nil {
		assert(tq.back == tq.front, "2")
		assert(tq.back.prev == nil, "3")
		return
	}
	nodes := make([]*tqItem[K, V], 0)
	for i := tq.front; i != nil; i = i.next {
		nodes = append(nodes, i)
	}
	index := 1
	for i := tq.back; i != nil; i = i.prev {
		assert(nodes[len(nodes)-index] == i, "4")
		index += 1
	}
}
