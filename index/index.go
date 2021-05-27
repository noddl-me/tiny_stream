package main

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"stream_test/stream"
)

func isNil(i interface{}) bool {
	return i == nil
}

type Stream struct {
	listeners []func(item interface{})
	Done      chan struct{}
	init      func(next func(item interface{}))
}

type ConcurrencyQueue struct {
	rw   sync.RWMutex
	data []interface{}
}

func (q *ConcurrencyQueue) Push(item interface{}) {
	q.rw.Lock()
	defer q.rw.Unlock()
	q.data = append(q.data, item)
}

func (q *ConcurrencyQueue) Pop() interface{} {
	q.rw.Lock()
	defer q.rw.Unlock()
	if len(q.data) == 0 {
		return nil
	}
	item := q.data[0]
	q.data = q.data[1:]
	return item
}

const (
	STOP      = 0
	RUNNING   = 1
	COMPLETED = 2
)

type TaskQueue struct {
	queue      ConcurrencyQueue
	status     int32
	f          func(item interface{})
	onComplete func()
}

func (t *TaskQueue) Push(item interface{}) {
	t.queue.Push(item)
}

func (t *TaskQueue) Run() {
	v := atomic.LoadInt32(&t.status)
	if v == RUNNING {
		return
	}

	if v == COMPLETED {
		panic("????")
	}
	atomic.CompareAndSwapInt32(&t.status, STOP, RUNNING)
	go func() {
		for {
			item := t.queue.Pop()
			if isNil(item) {
				atomic.CompareAndSwapInt32(&t.status, RUNNING, STOP)
				break
			}
			t.f(item)
		}

		if atomic.LoadInt32(&t.status) == COMPLETED {
			t.onComplete()
		}
	}()
}

func (t *TaskQueue) Complete() {
	status := atomic.LoadInt32(&t.status)
	if status == STOP {
		t.onComplete()
	} else {
		atomic.StoreInt32(&t.status, COMPLETED)
	}
}

func newStream() Stream {
	return Stream{
		Done: make(chan struct{}),
	}
}

func (s *Stream) ConcatMap(f func(item interface{}) *Stream) *Stream {
	newStream := newStream()
	queue := TaskQueue{
		f: func(item interface{}) {
			subStream := f(item)
			subStream.listeners = append(subStream.listeners, func(item interface{}) {
				newStream.Next(item)
			})
			subStream.init(subStream.Next)
			<-subStream.Done
		},
		onComplete: func() {
			newStream.Complete()
		},
	}

	s.listeners = append(s.listeners, func(item interface{}) {
		queue.Push(item)
		queue.Run()
	})

	go func() {
		<-s.Done
		queue.Complete()
	}()

	return &newStream
}

func (s *Stream) Complete() {
	s.Done <- struct{}{}
}

func (s *Stream) Next(item interface{}) {
	for _, f := range s.listeners {
		f(item)
	}
}

func modify(k []int) []int {
	for i := 0; i < 100000; i++ {
		k = append(k, 1, 2, 3, 4, 5)
	}
	fmt.Println((*reflect.SliceHeader)(unsafe.Pointer(&k)))
	return k
}

func main() {
	s := stream.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)
	s = s.
		Slice(2, 6).
		Filter(func(e stream.T) bool {
			return e.(int)%2 == 0
		}).
		Sort(func(left stream.T, right stream.T) int {
			return right.(int) - left.(int)
		})
	s.ForEach(func(e stream.T) {
		fmt.Println(e)
	})
	fmt.Println(s.Join(","))
}

func main1() {
	s := newStream()
	ss := s.ConcatMap(func(item interface{}) *Stream {
		inner := newStream()
		inner.init = func(next func(item interface{})) {
			wg := sync.WaitGroup{}
			for i := 0; i < 3; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					time.Sleep(time.Second * 2)
					next(item.(int) + 1)
				}(i)
			}
			go func() {
				wg.Wait()
				inner.Complete()
			}()
		}
		return &inner
	})

	ss.listeners = append(ss.listeners, func(item interface{}) {
		fmt.Printf("%v\n", item)
	})

	s.Next(1)
	s.Next(10)
	s.Next(100)
	s.Complete()

	defer func() {
		<-ss.Done
	}()
}
