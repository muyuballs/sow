package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	ver, err := ReadByte(conn)
	if err != nil {
		return
	}
	if ver == socksVer5 {
		target, err := handleSocks5(conn)
		if err != nil {
			log.Println(err)
			return
		}
		transferOverWebSocket(target, conn)
	} else if ver == socksVer4 {
		target, err := handleSocks4(conn)
		if err != nil {
			log.Println(err)
			return
		}
		transferOverWebSocket(target, conn)
	} else {
		conn.Close()
		log.Println("version", ver, "not supported")
	}
}

func main() {
	ln, err := net.Listen("tcp", ":40000")
	if err != nil {
		panic(err)
		return
	}
	log.Println("Start @", ln.Addr())
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		log.Println("new client:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
