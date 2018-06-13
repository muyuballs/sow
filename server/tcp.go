package main

import (
	slog "log"
	"net"
)

func handleTCP(c *Config) (err error) {
	log := slog.New(c.LogOut, "TCP ", LOG_FLAGS)
	laddr, err := net.ResolveTCPAddr("tcp", c.Listen)
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return
	}
	log.Printf("server start tcp://%v\n", laddr)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		err = conn.SetNoDelay(c.NoDelay)
		if err != nil {
			log.Println(err)
		}
		go transfer(conn, c, conn.RemoteAddr())
	}
	return
}
