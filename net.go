package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func runNet() {
	fmt.Println("net running on\t\t3000")
	l, err := net.Listen("tcp4", "127.0.0.1:3000")
	if err != nil {
		log.Fatalln("failed to listen 3000", err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			if err != io.EOF {
				log.Fatalln(err)
			}
		}
		go serveConn(c)
	}
}

var (
	part1 = []byte("HTTP/1.1 200 Ok\r\nContent-Type: text/plain; charset=utf-8\r\nDate: ")
	part2 = []byte("\r\nServer: net_\r\nContent-Length: 13\r\n\r\nHello, World!")
)

func serveConn(c net.Conn) {
	var (
		br  = bufio.NewReaderSize(c, 1024)
		bw  = bufio.NewWriterSize(c, 1024)
		err error
		buf []byte
		n   int
	)
	for {
		if n > 0 {
			_, err = br.Discard(n)
			if err == io.EOF {
				continue
			}
		} else {
			for {
				buf, _, err = br.ReadLine()
				if err != nil && err != io.EOF {
					break
				}
				n += len(buf) + 2
				// blank line -> len(buf) is 0
				if len(buf) == 0 {
					break
				}
				// request not ready, yet
			}
		}

		if err != nil {
			break
		}

		out := append(part1, ServerDate.Load().([]byte)...)
		out = append(out, part2...)

		_, err = bw.Write(out)
		if err != nil {
			break
		}

		if br == nil || br.Buffered() == 0 {
			if err = bw.Flush(); err != nil {
				break
			}
		}
	}

	c.Close()
}
