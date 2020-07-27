package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
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

		n := 1
		for {
			buf, err = br.Peek(n)
			if err != nil {
				break
			}
			if i := bytes.Index(buf, delimiter); i != -1 {
				_, err = br.Discard(br.Buffered())
				break
			}
			// request not ready, yet
			n += 1
		}

		if err != nil && err != io.EOF {
			log.Println("read error", err)
			break
		}

		if bw == nil {
			bw = acquireWriter(c)
		}

		out := append([]byte("HTTP/1.1 200 Ok\r\nContent-Type: text/plain; charset=utf-8\r\nDate: "), ServerDate.Load().([]byte)...)
		out = append(out, []byte("\r\nServer: net\r\nContent-Length: 13\r\n\r\nHello, World!")...)

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
