package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

var client *http.Client

type proxy struct {
	config proxyConfig
}

func newPorxy(c proxyConfig) *proxy {
	return &proxy{
		config: c,
	}
}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	resps := make([]*http.Response, len(p.config.EndPoints))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err: %s\n", err.Error())
	}

	defer r.Body.Close()

	var wg sync.WaitGroup

	for i, endpoint := range p.config.EndPoints {
		wg.Add(1)

		fullURL := endpoint + r.RequestURI
		index := i

		bb := &bytes.Buffer{}
		bb.Write(body)

		go func() {
			defer wg.Done()
			req, err := http.NewRequest(r.Method, fullURL, bb)
			for name, value := range r.Header {
				req.Header.Set(name, value[0])
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("client do err: %s, err: %s\n", fullURL, err.Error())
				return
			}

			resps[index] = resp
		}()
	}

	wg.Wait()

	// multiple response, choose the first not-nil resp to client
	wrReady := false
	for _, resp := range resps {
		if resp != nil {
			conn := &httpConnection{r, resp}
			if resp.StatusCode != 200 {
				printHTTP(conn)
			}

			if wrReady == true {
				io.Copy(ioutil.Discard, resp.Body)
				defer resp.Body.Close()

				continue
			}

			for k, v := range resp.Header {
				if len(v) != 0 {
					wr.Header().Set(k, v[0])
				}
			}

			wr.WriteHeader(resp.StatusCode)
			io.Copy(wr, resp.Body)
			defer resp.Body.Close()
			wrReady = true
		}
	}

}
