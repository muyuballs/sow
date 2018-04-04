package main

import (
	"io"
	"log"
	"time"

	kcp "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
)

func handleMux(conn io.ReadWriteCloser, c *Config) {
	// stream multiplex
	smuxConfig := smux.DefaultConfig()
	smuxConfig.MaxReceiveBuffer = 4194304
	smuxConfig.KeepAliveInterval = time.Duration(10) * time.Second

	mux, err := smux.Server(conn, smuxConfig)
	if err != nil {
		log.Println(err)
		return
	}
	defer mux.Close()
	for {
		p1, err := mux.AcceptStream()
		if err != nil {
			log.Println(err)
			return
		}
		go transfer(p1, c)
	}
}

func handleUDT(c *Config) (err error) {
	block, _ := kcp.NewNoneBlockCrypt(nil)
	lis, err := kcp.ListenWithOptions(c.Listen, block, 10, 3)
	if err != nil {
		return
	}
	log.Printf("server start udt://%v\n", lis.Addr())
	if err := lis.SetReadBuffer(c.SockBuf); err != nil {
		log.Println("SetReadBuffer:", err)
	}
	if err := lis.SetWriteBuffer(c.SockBuf); err != nil {
		log.Println("SetWriteBuffer:", err)
	}
	for {
		if conn, err := lis.AcceptKCP(); err == nil {
			log.Println("remote address:", conn.RemoteAddr())
			conn.SetStreamMode(true)
			conn.SetWriteDelay(true)
			conn.SetNoDelay(1, 20, 2, 1)
			conn.SetMtu(c.Mtu)
			conn.SetWindowSize(c.Swnd, c.Rwnd)
			conn.SetACKNoDelay(c.NoDelay)
			go handleMux(conn, c)
		} else {
			log.Printf("%+v", err)
		}
	}
}
