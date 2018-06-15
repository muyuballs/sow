package main

import (
	"fmt"
	"io"
	slog "log"
	"time"

	"github.com/muyuballs/sow/core"
)

func transfer2(target, origin string, c *core.Config, sconn io.ReadWriter, cconn io.ReadWriter) (err error) {
	sid := fmt.Sprintf("S%v ", time.Now().UTC().UnixNano())
	log := slog.New(c.LogOut, sid+" ", c.LOG_FLAGS)
	key := []byte(fmt.Sprintf("SOW-%v", time.Now().UTC().Format("200601021504")))
	log.Println("origin", origin)
	log.Println("target", target)
	rWriter := core.NewAESWriter(key, sconn)
	_, err = core.WriteLd(rWriter, []byte(c.Key))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = core.WriteLd(rWriter, []byte(target))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = core.WriteLd(rWriter, []byte(origin))
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		rReader := core.NewAESReader(key, sconn)
		n, e := io.Copy(cconn, rReader)
		log.Println("S->C", n, e)
	}()
	n, e := io.Copy(rWriter, cconn)
	log.Println("C->S", n, e)
	return
}
