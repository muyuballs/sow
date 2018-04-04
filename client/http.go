package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var (
	SowPpHeader = http.CanonicalHeaderKey("--SOW-PP--")
	transport   = &http.Transport{
		Proxy: resoloveProxy,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	np = make([]*regexp.Regexp, 0)
)

func resoloveProxy(req *http.Request) (*url.URL, error) {
	//	for _, reg := range np {
	//		if reg.MatchString(req.Host) {
	//			log.Println(req.Host, "In NP list")
	//			return nil, nil
	//		}
	//	}
	if pb, ok := req.Header[SowPpHeader]; ok {
		req.Header.Del(SowPpHeader)
		pp, err := url.Parse(pb[0])
		if err == nil {
			if pp.Host == "" || pp.Host == "0.0.0.0" {
				pp.Host = "127.0.0.1"
			}
		}
		return pp, err
	} else {
		log.Println(SowPpHeader, "not found")
	}
	return nil, nil
}
func handleHttp(preRead byte, conn *net.TCPConn, c *Config) {
	defer conn.Close()
	rd := bufio.NewReader(conn)
	req, err := http.ReadRequest(rd)
	if err != nil {
		log.Println(err)
		return
	}
	laddr, _ := conn.LocalAddr().(*net.TCPAddr)
	method := make([]byte, 0)
	method = append(method, preRead)
	method = append(method, []byte(req.Method)...)
	req.Method = string(method)
	req.Header.Add(SowPpHeader, fmt.Sprintf("socks5://%s", laddr.String()))
	bw := bufio.NewWriter(conn)
	defer bw.Flush()
	log.Println(req.Method, req.RequestURI)
	if req.Method == "CONNECT" {
		bw.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n")
		bw.Flush()
		log.Println(req.RequestURI)
		transferFn(req.RequestURI, conn, c)
	} else {
		resp, err := transport.RoundTrip(req)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		fmt.Fprintf(bw, "%d %s\r\n", resp.Status, resp.Proto)
		for k, v := range resp.Header {
			for _, vv := range v {
				fmt.Fprintf(bw, "%s: %s\r\n", k, vv)
			}
		}
		for _, c := range resp.Cookies() {
			fmt.Fprintf(bw, "Set-Cookie: %s\r\n", c.Raw)
		}
		fmt.Fprint(bw, "\r\n")
		_, err = io.Copy(bw, resp.Body)
		if err != nil && err != io.EOF {
			log.Println(err)
		}
	}
}
