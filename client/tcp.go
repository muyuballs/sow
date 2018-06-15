package main

import (
	slog "log"
	"net"

	"github.com/muyuballs/sow/core"
)

func transferByTCP(target string, conn *net.TCPConn, c *core.Config) {
	log := slog.New(c.LogOut, "TCP ", c.LOG_FLAGS)
	defer log.Println("job done")
	defer conn.Close()
	log.Println("dial server", c.Server)
	rconn, err := net.Dial("tcp", c.Server)
	if err != nil {
		log.Println(err)
		return
	}
	defer rconn.Close()
	transfer2(target, conn.RemoteAddr().String(), c, rconn, conn)
}
