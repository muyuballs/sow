package main

import (
	"log"
	"net"
)

func transferByTCP(target string, conn *net.TCPConn, c *Config) {
	defer log.Println("job done")
	defer conn.Close()
	log.Println("dial server")
	rconn, err := net.Dial("tcp", c.Server)
	if err != nil {
		log.Println(err)
		return
	}
	defer rconn.Close()
	transfer2(target, c.Key, rconn, conn, c.Zlib)
}
