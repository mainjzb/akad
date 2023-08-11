package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
)

// https://cdn-fastly.obsproject.com/downloads/OBS-Studio-29.1.3-Full-Installer-x64.exe
// https://client-updates-cdn77.badlion.net/Badlion%20Client%20Setup%203.16.0.exe
func main() {
	t, err := NewTask("https://cdn-fastly.obsproject.com/downloads/OBS-Studio-29.1.3-Full-Installer-x64.exe")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Run()
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
		return fmt.Errorf("status code err: %d", resp.StatusCode)
	}

	length := GetContentLength(resp.Header)
	bs := make([]byte, 0, length)
	resBytes := bytes.NewBuffer(bs)

	ctx, _ := context.WithCancel(context.Background())
	mr := &MyReader{
		ctx: ctx,
		r:   resp.Body,
	}

	resBytes.ReadFrom(mr)

	if resp.Header.Get("Accept-Ranges") == "bytes" {

	}

	return nil
}
