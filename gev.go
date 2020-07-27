package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/connection"
	"github.com/Allenxuxu/ringbuffer"
)

type gevServer struct{}

func (s *gevServer) OnConnect(c *connection.Connection) {
}
func (s *gevServer) OnMessage(c *connection.Connection, ctx interface{}, data []byte) (out []byte) {
	out = append(out, "HTTP/1.1 200 OK\r\nContent-Type: text/plain; charset=utf-8\r\nDate: "...)
	out = append(out, ServerDate.Load().([]byte)...)
	out = append(out, "\r\nServer: gev_\r\nContent-Length: 13\r\n\r\nHello, World!"...)
	return
}

func (s *gevServer) OnClose(c *connection.Connection) {
	//log.Println("OnClose")
}

func runGev() {
	fmt.Println("gev running on\t\t3002")

	s, err := gev.NewServer(new(gevServer),
		gev.Address("127.0.0.1:3002"),
		gev.Protocol(&httpProtocol{}),
	)

	if err != nil {
		log.Fatalln("failed to listen 3002", err)
	}

	s.Start()
}

type httpProtocol struct{}

func (d *httpProtocol) UnPacket(c *connection.Connection, buffer *ringbuffer.RingBuffer) (interface{}, []byte) {
	if i := bytes.Index(buffer.Bytes(), delimiter); i != -1 {
		buf := buffer.Bytes()
		buffer.RetrieveAll()
		return nil, buf
	}
	// request not ready, yet
	return nil, nil
}

func (d *httpProtocol) Packet(c *connection.Connection, data []byte) []byte {
	return data
}
