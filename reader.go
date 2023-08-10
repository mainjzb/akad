package main

import (
	"context"
	"io"

	log "github.com/sirupsen/logrus"
)

type MyReader struct {
	ctx context.Context
	r   io.ReadCloser
	cur int64 // current download size
	max int64 // max download size
}

func (r *MyReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		n, err = r.r.Read(p)

		r.cur += int64(n)
		log.Error("current download size: ", r.cur, " max download size: ", r.max)
		if r.max != 0 && r.cur >= r.max {
			return 0, io.EOF
		}

		return
	}
}

func (r *MyReader) Close() error {
	return r.r.Close()
}
