package main

import (
	"log"
	"net"
	"time"

	kcp "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
)

var dialUdt = func(server string) (*smux.Session, error) {
	block, _ := kcp.NewNoneBlockCrypt(nil)
	kcpconn, err := kcp.DialWithOptions(server, block, 10, 3)
	if err != nil {
		return nil, err
	}
	kcpconn.SetStreamMode(true)
	kcpconn.SetWriteDelay(true)
	kcpconn.SetNoDelay(1, 20, 2, 1)
	kcpconn.SetWindowSize(128, 512)
	kcpconn.SetMtu(1350)
	kcpconn.SetACKNoDelay(false)

	//	if err := kcpconn.SetDSCP(0); err != nil {
	//		log.Println("SetDSCP:", err)
	//	}
	if err := kcpconn.SetReadBuffer(4194304); err != nil {
		log.Println("SetReadBuffer:", err)
	}
	if err := kcpconn.SetWriteBuffer(4194304); err != nil {
		log.Println("SetWriteBuffer:", err)
	}
	smuxConfig := smux.DefaultConfig()
	smuxConfig.MaxReceiveBuffer = 4194304
	smuxConfig.KeepAliveInterval = time.Duration(10) * time.Second
	session, err := smux.Client(kcpconn, smuxConfig)
	return session, err
}

func transferByUDT(target string, conn *net.TCPConn, c *Config) {
	log.Println("stream opened")
	defer log.Println("stream closed")
	defer conn.Close()
	session, err := dialUdt(c.Server)
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
