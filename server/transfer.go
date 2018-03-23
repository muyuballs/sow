package main

import (
	"compress/flate"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/muyuballs/sow/crypt"
)

func readld(r io.Reader) (dat []byte, err error) {
	buf := make([]byte, 4)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		log.Println(err)
		return
	}
	tl := binary.BigEndian.Uint32(buf)
	if tl > 1024 {
		log.Println("seg too long", tl)
		return
	}
	dat = make([]byte, int(tl))
	_, err = io.ReadFull(r, dat)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func checkKey(r io.Reader, xkey string) error {
	key, err := readld(r)
	if err != nil {
		return err
	}
	if xkey != string(key) {
		return fmt.Errorf("key [%v] is invalid", string(key))
	}
	return nil
}

func transfer(conn io.ReadWriteCloser, c *Config) (err error) {
	defer log.Println("job done")
	defer conn.Close()
	key := []byte(fmt.Sprintf("SOW-%v", time.Now().UTC().Format("200601021504")))
	var rReader io.Reader
	var rWriter io.Writer
	if c.Zlib {
		rReader = crypt.NewAESReader(key, conn)
		rReader, err = zlib.NewReader(rReader)
		if err != nil {
			log.Println(err)
			return
		}
		rWriter = crypt.NewAESWriter(key, conn)
		rWriter, err = zlib.NewWriterLevel(rWriter, flate.BestCompression)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		rReader = crypt.NewAESReader(key, conn)
		rWriter = crypt.NewAESWriter(key, conn)
	}
	err = checkKey(rReader, c.Key)
	if err != nil {
		log.Println(err)
		return
	}
	buf, err := readld(rReader)
	if err != nil {
		log.Println(err)
		return
	}
	target := string(buf)
	log.Println(target)
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
	go io.Copy(rconn, rReader)
	buf = make([]byte, 32*1024)
	for {
		nr, er := rconn.Read(buf)
		if nr > 0 {
			nw, ew := rWriter.Write(buf[0:nr])
			if ew != nil {
				log.Println(ew)
				break
			}
			if nr != nw {
				log.Println("short write")
				break
			}
			if c.Zlib {
				fe := rWriter.(*zlib.Writer).Flush()
				if fe != nil {
					log.Println(fe)
					break
				}
			}
		}
		if er != nil {
			log.Println(er)
			break
		}
	}
	return
}
