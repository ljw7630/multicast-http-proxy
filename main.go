package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

var gCurrentPath string

type httpConnection struct {
	Request  *http.Request
	Response *http.Response
}

func printHTTP(conn *httpConnection) {
	fmt.Printf("%v %v\n", conn.Request.Method, conn.Request.RequestURI)
	for k, v := range conn.Request.Header {
		fmt.Printf("%s : %s\n", k, v)
	}
	fmt.Println("------------------------------")
	fmt.Printf("HTTP/1.1 %v\n", conn.Response.Status)
	for k, v := range conn.Response.Header {
		fmt.Printf("%s : %s\n", k, v)
	}
	fmt.Println("==============================")
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "conf", "", "config file path")

	flag.Parse()

	c, err := newProxyConfig(configPath)
	if err != nil {
		fmt.Printf("new proxy config err: %s", err.Error())
		return
	}

	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(10) * time.Second,
	}

	p := newPorxy(*c)
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.BindPort), p)
	if err != nil {
		fmt.Printf("ListenAndServe: %s\n", err.Error())
	}
}
