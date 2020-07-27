package main

import (
	"bytes"
	"fmt"
	"log"

	gnet "github.com/panjf2000/gnet"
)

func runGnet() {
	hc := &httpCodec{delimiter: delimiter}

	// Start serving!
	fmt.Println("gnet running on\t\t3001")
	log.Fatalln(gnet.Serve(new(gnetServer), "tcp://127.0.0.1:3001",
		gnet.WithMulticore(true),
		gnet.WithCodec(hc),
	))
}

type gnetServer struct {
	*gnet.EventServer
}

type httpCodec struct {
	delimiter []byte
}

func (hc *httpCodec) Encode(c gnet.Conn, buf []byte) (out []byte, err error) {
	return buf, nil
}

func (hc *httpCodec) Decode(c gnet.Conn) (out []byte, err error) {
	buf := c.Read()
	if buf == nil {
		return
	}
	c.ResetBuffer()

	// process the pipeline
	var i int
pipeline:
	if i = bytes.Index(buf, hc.delimiter); i != -1 {
		out = append(out, "HTTP/1.1 200 OK\r\nContent-Type: text/plain; charset=utf-8\r\nDate: "...)
		out = append(out, ServerDate.Load().([]byte)...)
		out = append(out, "\r\nServer: gnet\r\nContent-Length: 13\r\n\r\nHello, World!"...)
		buf = buf[i+4:]
		goto pipeline
	}
	// request not ready, yet
	return
}

func (hs *gnetServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	// handle the request
	out = frame
	return
}
