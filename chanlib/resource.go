package chanlib

import (
	"log"
	"net/http"
	"time"
)

type Resource struct {
	url      string
	errCount int
}

func NewResource(url string) *Resource {
	return &Resource{url: url, errCount: 0}
}

func (r *Resource) GetUrl() string {
	return r.url
}

func (r *Resource) Poll() string {
	res, err := http.Head(r.url)
	if err != nil {
		r.errCount++
		log.Println("")
		log.Println("")
		log.Println("Error")
		log.Println(r.url, err)
		return err.Error()
	}
	r.errCount = 0
	log.Println("")
	log.Println("")
	log.Println("Info:")
	log.Println(r.url, " reques header completed")
	return res.Status
}

func (r *Resource) Sleep(baseInterval time.Duration, done chan<- *Resource) {
	var errTime = baseInterval / 10
	time.Sleep(baseInterval + errTime*time.Duration(r.errCount))
	done <- r
}
