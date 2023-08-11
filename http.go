package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func GetContentLength(header http.Header) int64 {
	v := header.Get("Content-Length")
	if v == "" {
		return 0
	}
	length, _ := strconv.Atoi(v)
	return int64(length)
}

func Connect(u url.URL, ip string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	host := u.Host
	// todo 考虑host中有端口的情况
	u.Host = ip

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Host = host
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code err: %d", resp.StatusCode)
	}
	return resp, nil
}

func ConnectWithRange(u url.URL, ip string, start, end int64) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	host := u.Host
	// todo 考虑host中有端口的情况
	u.Host = ip

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Host = host
	req.Header.Set("Range", Range(start, end))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status code err: %d", resp.StatusCode)
	}
	return resp, nil
}

func Range(start, end int64) string {
	if end == 0 {
		return fmt.Sprintf("bytes=%d-", start)
	}
	return fmt.Sprintf("bytes=%d-%d", start, end-1)
}
