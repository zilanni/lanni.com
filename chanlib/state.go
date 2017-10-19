package chanlib

import (
	"log"
	"time"
)

type State struct {
	Url    string
	Status string
}

func StateMonitor(updateInternal time.Duration) chan<- State {
	stateMap := make(map[string]string)
	stateChan := make(chan State, 2)
	timer := time.NewTicker(updateInternal)
	go func() {
		for {
			select {
			case <-timer.C:
				logState(stateMap)
			case s := <-stateChan:
				stateMap[s.Url] = s.Status
			}
		}
	}()
	return stateChan
}

// logState prints a state map.

// logState 打印出一个状态映射。
func logState(s map[string]string) {
	log.Println("Current state:")
	for k, v := range s {
		log.Printf(" %s %s", k, v)
	}
}
