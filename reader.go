package main

import (
	"context"
	"io"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

type MyReader struct {
	ctx context.Context
	r   io.ReadCloser
	cur int64 // current download size
	// Max int64 // Max download size
	Max   atomic.Int64 // Max download size
	Start int64
}

func NewMyReader(r io.ReadCloser, start, end int64) *MyReader {
	m := &MyReader{
		ctx:   context.Background(),
		r:     r,
		Start: start,
	}
	m.Max.Store(end - start)

	return m
}

func (r *MyReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		n, err = r.r.Read(p)

		r.cur += int64(n)
		log.Error("current download size: ", r.cur, " max download size: ", r.Max)
		if r.Max.Load() != 0 && r.cur >= r.Max.Load() {
			return 0, io.EOF
		}

		return
	}
}

func (r *MyReader) Close() error {
	return r.r.Close()
}
