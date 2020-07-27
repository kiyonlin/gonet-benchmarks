package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var (
	part1 = []byte("HTTP/1.1 200 Ok\r\nContent-Type: text/plain; charset=utf-8\r\nDate: ")
	part2 = []byte("\r\nServer: net_\r\nContent-Length: 13\r\n\r\nHello, World!")
)

var readerPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReaderSize(nil, 4096)
	},
}
var writerPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewWriterSize(nil, 4096)
	},
}

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

func serveConn(c net.Conn) {
	br := acquireReader(c)
	var (
		bw  *bufio.Writer
		err error
		buf []byte
	)
	for {
		buf, err = br.Peek(4)
		if len(buf) == 0 {
			if err == io.EOF {
				continue
			}
			break
		}

		for {
			buf, _, err = br.ReadLine()
			if err != nil {
				break
			}
			// blank line -> len(buf) is 0
			if len(buf) == 0 {
				// Discard all read data
				_, err = br.Discard(br.Buffered())
				break
			}
			// request not ready, yet
		}

		if err != nil && err != io.EOF {
			log.Println("read error", err)
			break
		}

		if bw == nil {
			bw = acquireWriter(c)
		}

		out := append(part1, ServerDate.Load().([]byte)...)
		out = append(out, part2...)

		_, err = bw.Write(out)
		if err != nil {
			log.Println("write error", err)
			break
		}

		if br == nil || br.Buffered() == 0 {
			if err = bw.Flush(); err != nil {
				break
			}
		}
	}

	releaseReader(br)

	if bw != nil {
		releaseWriter(bw)
	}

	c.Close()
}

func acquireReader(c net.Conn) *bufio.Reader {
	r := readerPool.Get().(*bufio.Reader)
	r.Reset(c)
	return r
}

func releaseReader(r *bufio.Reader) {
	readerPool.Put(r)
}

func acquireWriter(c net.Conn) *bufio.Writer {
	w := writerPool.Get().(*bufio.Writer)
	w.Reset(c)
	return w
}

func releaseWriter(w *bufio.Writer) {
	writerPool.Put(w)
}
