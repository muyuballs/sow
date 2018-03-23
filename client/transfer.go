package main

import (
	"compress/flate"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/muyuballs/sow/crypt"
)

func writeld(w io.Writer, b []byte) (n int, err error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(b)))
	n, err = w.Write(buf)
	if err != nil {
		log.Println(err)
		return
	}
	n, err = w.Write([]byte(b))
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func transfer2(target, secret string, sconn io.ReadWriter, cconn io.ReadWriter, compress bool) (err error) {
	key := []byte(fmt.Sprintf("SOW-%v", time.Now().UTC().Format("200601021504")))
	log.Println("create writer")
	var rWriter io.Writer
	if compress {
		rWriter = crypt.NewAESWriter(key, sconn)
		rWriter, err = zlib.NewWriterLevel(rWriter, flate.BestCompression)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		rWriter = crypt.NewAESWriter(key, sconn)
	}
	_, err = writeld(rWriter, []byte(secret))
	if err != nil {
		log.Println(err)
		return
	}
	_, err = writeld(rWriter, []byte(target))
	if err != nil {
		log.Println(err)
		return
	}
	if compress {
		fe := rWriter.(*zlib.Writer).Flush()
		if fe != nil {
			log.Println(fe)
		}
	}
	go func() {
		log.Println("create reader")
		var rReader io.Reader
		if compress {
			rReader = crypt.NewAESReader(key, sconn)
			rReader, err = zlib.NewReader(rReader)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			rReader = crypt.NewAESReader(key, sconn)
		}
		log.Println("start client <- server")
		io.Copy(cconn, rReader)
	}()
	log.Println("start client -> server")
	buf := make([]byte, 32*1024)
	for {
		nr, er := cconn.Read(buf)
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
			if compress {
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
