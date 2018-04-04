package main

import (
	"log"
	"net"
	"time"

	kcp "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
)

var dialUdt = func(c *Config) (*smux.Session, error) {
	block, _ := kcp.NewNoneBlockCrypt(nil)
	kcpconn, err := kcp.DialWithOptions(c.Server, block, 10, 3)
	if err != nil {
		return nil, err
	}
	kcpconn.SetStreamMode(true)
	kcpconn.SetWriteDelay(true)
	kcpconn.SetNoDelay(1, 20, 2, 1)
	kcpconn.SetWindowSize(c.Swnd, c.Rwnd)
	kcpconn.SetMtu(c.Mtu)
	kcpconn.SetACKNoDelay(false)

	if err := kcpconn.SetReadBuffer(c.SockBuf); err != nil {
		log.Println("SetReadBuffer:", err)
	}
	if err := kcpconn.SetWriteBuffer(c.SockBuf); err != nil {
		log.Println("SetWriteBuffer:", err)
	}
	smuxConfig := smux.DefaultConfig()
	smuxConfig.MaxReceiveBuffer = c.SockBuf
	smuxConfig.KeepAliveInterval = time.Duration(10) * time.Second
	session, err := smux.Client(kcpconn, smuxConfig)
	return session, err
}

func transferByUDT(target string, conn *net.TCPConn, c *Config) {
	log.Println("stream opened")
	defer log.Println("stream closed")
	defer conn.Close()
	session, err := dialUdt(c)
	if err != nil {
		log.Println(err)
		return
	}
	defer session.Close()
	p2, err := session.OpenStream()
	if err != nil {
		log.Println(err)
		return
	}
	defer p2.Close()
	transfer2(target, c.Key, p2, conn, c.Zlib)
}
