package main

import (
	"fmt"
	"io"
	slog "log"
	"net"
	"time"

	"github.com/muyuballs/sow/core"
)

func checkKey(r io.Reader, xkey string, log *slog.Logger) error {
	key, err := core.ReadLd(r)
	if err != nil {
		return err
	}
	if xkey != string(key) {
		return fmt.Errorf("key [%v] is invalid", string(key))
	}
	return nil
}

func transfer(conn io.ReadWriteCloser, c *core.Config, remoteAddr net.Addr) (err error) {
	sid := fmt.Sprintf("S%v ", time.Now().UTC().UnixNano())
	log := slog.New(c.LogOut, sid+" ", c.LOG_FLAGS)
	log.Println("client", remoteAddr)
	defer log.Println("job done")
	defer conn.Close()
	key := []byte(fmt.Sprintf("SOW-%v", time.Now().UTC().Format("200601021504")))
	rReader := core.NewAESReader(key, conn)
	rWriter := core.NewAESWriter(key, conn)
	err = checkKey(rReader, c.Key, log)
	if err != nil {
		log.Println(err)
		return
	}
	buf, err := core.ReadLd(rReader)
	if err != nil {
		log.Println(err)
		return
	}
	target := string(buf)
	log.Println("target", target)
	buf, err = core.ReadLd(rReader)
	if err != nil {
		log.Println(err)
		return
	}
	origin := string(buf)
	log.Println("origin", origin)
	rAddr, err := net.ResolveTCPAddr("tcp", target)
	if err != nil {
		log.Println(err)
		return
	}
	rconn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer rconn.Close()
	go func() {
		c, err := io.Copy(rconn, rReader)
		log.Println("C->D", c, err)
	}()
	buf = make([]byte, 32*1024)
	cc, err := io.Copy(rWriter, rconn)
	log.Println("D->C", cc, err)
	return
}
