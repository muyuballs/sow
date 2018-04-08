package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

func handMomoConnection(c *net.TCPConn, conf *Config) {
	defer c.Close()
	buf := make([]byte, 4)
	_, err := io.ReadFull(c, buf)
	if err != nil {
		log.Println(err)
		return
	}
	targetSize := binary.BigEndian.Uint32(buf)
	buf = make([]byte, int(targetSize))
	_, err = io.ReadFull(c, buf)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("momo target", string(buf))

	transferFn(string(buf), c, conf)
}
