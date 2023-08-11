package main

import (
	"bytes"
	"net/url"
	"sync"

	log "github.com/sirupsen/logrus"
)

const MAXTask = 8

type Task struct {
	url *url.URL
	l   []*MyReader
	wg  *sync.WaitGroup

	contentLength int64
}

func NewTask(u string) (*Task, error) {
	_url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	return &Task{
		wg:  &sync.WaitGroup{},
		url: _url,
	}, nil
}

func (t *Task) Run() {
	resp, err := Connect(*t.url, "151.101.1.91")
	if err != nil {
		log.Error("DownloadIP err: ", err.Error())
		// todo
		return
	}
	// only for first
	t.contentLength = GetContentLength(resp.Header)

	t.l = append(t.l, NewMyReader(resp.Body, 0, t.contentLength))

	resBytes := bytes.NewBuffer(nil)
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		resBytes.ReadFrom(t.l[0])
	}()

	// second task
	// second task
	// second task

	t.DoubleTask()
	t.DoubleTask()
	t.DoubleTask()
	t.DoubleTask()
	t.DoubleTask()

	t.wg.Wait()
}

func (t *Task) SpawnChild(parent *MyReader) (*MyReader, error) {
	start := parent.Start + parent.Max.Load()/2
	end := parent.Start + parent.Max.Load()

	resp, err := ConnectWithRange(*t.url, "151.101.65.91", start, end)
	if err != nil {
		return nil, err
	}

	parent.Max.Store(parent.Max.Load() / 2)

	reader := NewMyReader(resp.Body, start, end)

	resBytes2 := bytes.NewBuffer(nil)
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		resBytes2.ReadFrom(t.l[1])
	}()

	return reader, nil
}

func (t *Task) DoubleTask() {
	res := make([]*MyReader, 0, len(t.l)*2)
	for _, l := range t.l {
		res = append(res, l)
		child, err := t.SpawnChild(l)
		if err != nil {
			log.Error(err)
			// todo 重新创建child
			continue
		}

		res = append(res, child)
	}
	t.l = res
}
