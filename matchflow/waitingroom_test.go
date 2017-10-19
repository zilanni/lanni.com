package matchflow

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

type Waiter struct {
	Matched     []Player
	ApplyFail   []Player
	WaitTimeout []Player
	mlock       chan bool
	alock       chan bool
	wlock       chan bool
	mCount      int64
	aCount      int64
	wCount      int64
}

func NewWaiter() *Waiter {
	return &Waiter{Matched: make([]Player, 0, 100), ApplyFail: make([]Player, 0, 10), WaitTimeout: make([]Player, 0, 50), mlock: make(chan bool, 1), alock: make(chan bool, 1), wlock: make(chan bool, 1)}
}

func (w *Waiter) ServePlayers(players []Player) {
	w.mlock <- true
	//w.Matched = append(w.Matched, players...)
	w.mCount += int64(len(players))
	<-w.mlock
}

func (w *Waiter) NoticePlayerWaitTimeout(player Player) {
	w.wlock <- true
	w.wCount++
	<-w.wlock
	//w.WaitTimeout = append(w.WaitTimeout, player)
}

func (w *Waiter) NoticePlayerApplyTimeout(player Player) {
	//w.ApplyFail = append(w.ApplyFail, player)
	w.alock <- true
	w.aCount++
	<-w.alock
}

func TestApplyIn(t *testing.T) {
	waiter := NewWaiter()
	wr := WaitingRoom{Waiter: waiter, PCount: 4, WTimeout: time.Millisecond * 50, ATimeout: time.Millisecond * 10}
	wr.Start()
	var counter int64 = 0
	for i := 1; i <= 100; i++ {
		go func(index int) {
			for j := 1; j <= 10000; j++ {
				id := 1000000*index + j
				player := Player{Id: int64(id)}
				wr.ApplyIn(player)
				atomic.AddInt64(&counter, 1)
				time.Sleep(time.Millisecond * time.Duration(index))
			}
		}(i)
	}
	time.Sleep(time.Second * 10)
	wr.Stop()
	totalCount := waiter.aCount + waiter.mCount + waiter.wCount
	if totalCount != counter {
		t.Errorf("%d ApplyIn want %d, get %d", 100*10000, counter, totalCount)
	}
	//t.Errorf("Reverse(%q) == %q, want %q, c.in, got, c.want")
}

func BenchmarkApplyIn(t *testing.B) {
	waiter := NewWaiter()
	wr := WaitingRoom{Waiter: waiter, PCount: 4, WTimeout: time.Millisecond * 50, ATimeout: time.Millisecond * 10}
	wr.Start()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		player := Player{Id: int64(i)}
		wr.ApplyIn(player)
	}
}

func BenchmarkApplyInPalla(t *testing.B) {
	waiter := NewWaiter()
	wr := WaitingRoom{Waiter: waiter, PCount: 4, WTimeout: time.Millisecond * 50, ATimeout: time.Millisecond * 10}
	wr.Start()
	t.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := rand.Int63()
			player := Player{Id: id}
			wr.ApplyIn(player)
		}
	})
}
