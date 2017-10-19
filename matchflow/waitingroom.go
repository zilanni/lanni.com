package matchflow

import (
	"time"
)

type Player struct {
	Name string
	Id   int64
}

type PlayerStatus struct {
	WPlayer   Player
	ApplyTime time.Time
}

type Iwaiter interface {
	ServePlayers(players []Player)
	NoticePlayerWaitTimeout(player Player)
	NoticePlayerApplyTimeout(player Player)
}

type WaitingRoom struct {
	Waiter       Iwaiter
	PCount       int
	WTimeout     time.Duration //Wait Timeout
	ATimeout     time.Duration //Apply Timeout
	requestQueue chan PlayerStatus
	outPlayers   chan PlayerStatus //outPlayers 匹配窗口，窗口大小比PCount小1
	ticker       *time.Ticker
	isworking    bool
	workingLock  chan bool
}

func (wr *WaitingRoom) ApplyIn(player Player) {
	if wr.isworking {
		ps := PlayerStatus{WPlayer: player, ApplyTime: time.Now()}
		select {
		case wr.requestQueue <- ps:
		case <-time.After(wr.ATimeout):
			wr.Waiter.NoticePlayerApplyTimeout(player)
		}
	} else {
		wr.Waiter.NoticePlayerApplyTimeout(player)
	}
}

func (wr *WaitingRoom) Start() {
	wr.requestQueue = make(chan PlayerStatus, 10000)
	wr.outPlayers = make(chan PlayerStatus, wr.PCount-1)
	wr.workingLock = make(chan bool, 1)
	if wr.ticker == nil {
		wr.ticker = time.NewTicker(wr.WTimeout / 3)
	}
	wr.isworking = true
	go func() {
		for wr.isworking {
			wr.workingLock <- true
			if wr.isworking {
				wr.poll()
			}
			<-wr.workingLock
		}
	}()
}

func (wr *WaitingRoom) Stop() {
	wr.workingLock <- true
	wr.isworking = false
	more := true
	for more {
		select {
		case ps := <-wr.requestQueue:
			wr.Waiter.NoticePlayerWaitTimeout(ps.WPlayer)
		default:
			more = false
		}
	}

	more = true
	for more {
		select {
		case ps := <-wr.outPlayers:
			wr.Waiter.NoticePlayerWaitTimeout(ps.WPlayer)
		default:
			more = false
		}
	}

	<-wr.workingLock
}

func (wr *WaitingRoom) poll() {
	select {
	case ps := <-wr.requestQueue:
		select {
		case wr.outPlayers <- ps:
		default:
			players := make([]Player, wr.PCount)
			for i := 0; i < wr.PCount-1; i++ {
				player := <-wr.outPlayers
				players[i] = player.WPlayer
			}
			players[wr.PCount-1] = ps.WPlayer
			wr.Waiter.ServePlayers(players)
		}
	case <-wr.ticker.C:
		count := len(wr.outPlayers)
		for i := 0; i < count; i++ {
			ps := <-wr.outPlayers
			if time.Now().After(ps.ApplyTime.Add(wr.WTimeout)) {
				wr.Waiter.NoticePlayerWaitTimeout(ps.WPlayer)
			} else {
				wr.outPlayers <- ps
			}
		}
	}
}
