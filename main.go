package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// https://cdn-fastly.obsproject.com/downloads/OBS-Studio-29.1.3-Full-Installer-x64.exe

func main() {
	u, err := url.Parse("https://cdn-fastly.obsproject.com/downloads/OBS-Studio-29.1.3-Full-Installer-x64.exe")
	fmt.Println(u, err)
	DownloadIP(u, "151.101.1.91")

}

func DownloadIP(u *url.URL, ip string) error {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	host := u.Host
	// todo 考虑host中有端口的情况
	u.Host = ip

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Host = host
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error("status code err: ", resp.StatusCode)
		return err
	}

	length := GetContentLength(resp.Header)
	bs := make([]byte, 0, length)
	resBytes := bytes.NewBuffer(bs)

	ctx, _ := context.WithCancel(context.Background())
	mr := &MyReader{
		ctx: ctx,
		r:   resp.Body,
	}

	go resBytes.ReadFrom(mr)

	vs := resp.Header.Values("Accept-Ranges")
	for _, v := range vs {
		if strings.ToLower(v) == "bytes" {
			// 启动子线程
		}
	}

	return nil
}

func GetContentLength(header http.Header) int {
	// todo chunck return -1
	v := header.Get("Content-Length")
	if v == "" {
		return 0
	}
	length, _ := strconv.Atoi(v)
	return length
}

type MyReader struct {
	ctx context.Context
	r   io.ReadCloser
	n   int64
	max int64 // max download size
}

func (r *MyReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		n, err = r.r.Read(p)
		r.n += int64(n)
		if r.max != 0 && r.n > r.max {
			return 0, io.EOF
		}
		return
	}
}

func (r *MyReader) Close() error {
	return r.r.Close()
}

func ReadFrom(b *bytes.Buffer, r io.Reader) (n int64, err error) {
	for {
		// tmp := make([]byte, 1024)
		// n, err := r.Read(tmp)
		// if n < 0 {
		// 	panic("err n < 0")
		// }
		//
		// b.Write(tmp[:n])
		//
		// b.Grow(1024)
		// b.Bytes()
		//
		// m, e := r.Read(b[b.Len():b.Cap()])
		//
		// b.buf = b.buf[:i+m]
		// n += int64(m)
		// if e == io.EOF {
		// 	return n, nil // e is EOF, so return nil explicitly
		// }
		// if e != nil {
		// 	return n, e
		// }
	}
}
