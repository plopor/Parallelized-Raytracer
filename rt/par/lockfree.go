package par

import (
	"sync"
	"sync/atomic"
)

type Task struct {
	X, Y int
}

type AtomicStampedReference struct {
	val   int64
	stamp int64
}

func NewAtomicStampedReference(initialVal, initialStamp int64) *AtomicStampedReference {
	return &AtomicStampedReference{
		val:   initialVal,
		stamp: initialStamp,
	}
}

func (asr *AtomicStampedReference) Get() (int64, int64) {
	return asr.val, asr.stamp
}

func (asr *AtomicStampedReference) CompareAndSwap(expectedVal, newVal, expectedStamp, newStamp int64) bool {
	if asr.val == expectedVal && asr.stamp == expectedStamp {
		atomic.StoreInt64((*int64)(&asr.val), int64(newVal))
		atomic.StoreInt64((*int64)(&asr.stamp), int64(newStamp))
		return true
	}
	return false
}

func (asr *AtomicStampedReference) set(newVal, newStamp int64) {
	atomic.StoreInt64((*int64)(&asr.val), int64(newVal))
	atomic.StoreInt64((*int64)(&asr.stamp), int64(newStamp))
}

type BDEQueue struct {
	top    AtomicStampedReference
	bottom atomic.Int32
	tasks  [48000000]Task
}

func NewBDEQueue(tasks []Task) *BDEQueue {
	q := &BDEQueue{
		top: AtomicStampedReference{val: 0, stamp: 0},
	}
	for _, task := range tasks {
		q.PushBottom(task)
	}
	return q
}

func (q *BDEQueue) PushBottom(task Task) {
	q.tasks[q.bottom.Load()] = task
	q.bottom.Add(1)
}

func (q *BDEQueue) PopTop() (Task, bool) {
	oldTop, oldStamp := q.top.Get()
	newTop := oldTop + 1
	newStamp := oldStamp + 1

	if int(q.bottom.Load()) <= int(oldTop) {
		return Task{}, false
	}
	task := q.tasks[oldTop]
	if q.top.CompareAndSwap(oldTop, newTop, oldStamp, newStamp) {
		return task, true
	}
	return Task{}, false
}

func (q *BDEQueue) PopBottom() (Task, bool) {
	if q.bottom.Load() == 0 {
		return Task{}, false
	}
	q.bottom.Add(-1)
	task := q.tasks[q.bottom.Load()]
	oldTop, oldStamp := q.top.Get()
	newTop := 0
	newStamp := oldStamp + 1
	if int(q.bottom.Load()) > int(oldTop) {
		return task, true
	}
	if int(q.bottom.Load()) == int(oldTop) {
		q.bottom.Store(0)
		if q.top.CompareAndSwap(oldTop, int64(newTop), oldStamp, newStamp) {
			return task, true
		}
	}
	q.top.set(int64(newTop), newStamp)
	q.bottom.Store(0)
	return Task{}, false
}

func Balance(self, victim *BDEQueue) {
	if self.Size() > victim.Size() {
		return
	}
	diff := victim.Size() - self.Size()
	if diff > victim.Size()/2 {
		for victim.Size() > self.Size() {
			value, _ := victim.PopTop()
			self.PushBottom(value)
		}
	}
}

func (q *BDEQueue) Size() int {
	bottom := q.bottom.Load()
	top, _ := q.top.Get()
	return int(bottom) - int(top)
}

type Barrier struct {
	Count int
	Total int
	mutex sync.Mutex
	cond  *sync.Cond
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{
		Total: n,
	}
	b.cond = sync.NewCond(&b.mutex)
	return b
}

func (b *Barrier) Signal() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.Count++
	if b.Count == b.Total {
		b.cond.Broadcast()
	}
}

func (b *Barrier) Wait() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for b.Count < b.Total {
		b.cond.Wait()
	}
}
