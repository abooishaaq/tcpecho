package main

import (
	"net"
	"sync"
)

type Conns struct {
	mutex sync.Mutex
	conns []*net.Conn
}

var conns Conns

func handleConn(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			break
		}
		str := string(buf[0 : n-1])
		for _, c := range conns.conns {
			(*c).Write([]byte(str + "\n"))
		}
	}

	for i, c := range conns.conns {
		if *c == conn {
			conns.mutex.Lock()
			conns.conns = append(conns.conns[:i], conns.conns[i+1:]...)
			conns.mutex.Unlock()
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":1337")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		conns.mutex.Lock()
		conns.conns = append(conns.conns, &conn)
		conns.mutex.Unlock()
		go handleConn(conn)
	}
}
