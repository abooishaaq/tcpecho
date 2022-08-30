package main

import "net"

var conns []*net.Conn

func handleConn(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			break
		}
		str := string(buf[0 : n-1])
		for _, c := range conns {
			(*c).Write([]byte(str + "\n"))
		}
	}

	for i, c := range conns {
		if *c == conn {
			conns = append(conns[:i], conns[i+1:]...)
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
		conns = append(conns, &conn)
		go handleConn(conn)
	}
}
