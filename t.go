package main

import (
	"bytes"
	"context"
	"net/url"
	"sync"

	log "github.com/sirupsen/logrus"
)

const MAXTask = 8

type Task struct {
	url *url.URL
	l   []*MyReader

	contentLength int64
}

func NewTask(u string) (*Task, error) {
	_url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return &Task{
		url: _url,
	}, nil
}

func (t *Task) Run() {
	wg := &sync.WaitGroup{}

	resp, err := Connect(*t.url, "143.244.51.206")
	if err != nil {
		log.Error("DownloadIP err: ", err.Error())
		// todo
		return
	}
	// only for first
	t.contentLength = GetContentLength(resp.Header)

	t.l = append(t.l, &MyReader{
		ctx: context.Background(),
		r:   resp.Body,
		max: t.contentLength / 2,
	})

	resBytes := bytes.NewBuffer(nil)
	wg.Add(1)
	go func() {
		defer wg.Done()
		resBytes.ReadFrom(t.l[0])
	}()

	// second task
	// second task
	// second task

	resp2, err := ConnectWithRange(*t.url, "89.187.187.11", t.contentLength/2+1, t.contentLength)
	if err != nil {
		log.Error("DownloadIP err: ", err.Error())
		// todo
		return
	}

	t.l = append(t.l, &MyReader{
		ctx: context.Background(),
		r:   resp2.Body,
		max: t.contentLength - t.contentLength/2,
	})

	resBytes2 := bytes.NewBuffer(nil)
	wg.Add(1)
	go func() {
		defer wg.Done()
		resBytes2.ReadFrom(t.l[1])
	}()

	wg.Wait()
}
