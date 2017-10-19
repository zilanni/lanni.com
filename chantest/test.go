package main

import (
	"time"

	"lanni.com/chanlib"
)

const (
	numPollers     = 2                // number of Poller goroutines to launch // Poller Go程的启动数
	pollInterval   = 60 * time.Second // how often to poll each URL            // 轮询每一个URL的频率
	statusInterval = 10 * time.Second // how often to log status to stdout     // 将状态记录到标准输出的频率
	//errTimeout     = 10 * time.Second // back-off timeout on error             // 回退超时的错误
)

var urls = []string{
	"http://golang.org/",
	"http://www.baidu.com/",
	"http://www.sina.com.cn/",
}

func Poll(in <-chan *chanlib.Resource, out chan<- *chanlib.Resource, states chan<- chanlib.State) {
	for r := range in {
		status := r.Poll()
		states <- chanlib.State{Url: r.GetUrl(), Status: status}
		out <- r
	}
}

func main() {
	request := make(chan *chanlib.Resource, 2)
	completed := make(chan *chanlib.Resource, 2)
	states := chanlib.StateMonitor(statusInterval)

	for i := 0; i < numPollers; i++ {
		go Poll(request, completed, states)
	}

	go func() {
		for _, url := range urls {
			request <- chanlib.NewResource(url)
		}
	}()

	for c := range completed {
		c.Sleep(pollInterval, request)
	}
}
