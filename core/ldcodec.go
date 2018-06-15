package core

import (
	"encoding/binary"
	"io"
)

func WriteLd(w io.Writer, b []byte) (n int, err error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(len(b)))
	n, err = w.Write(buf)
	if err != nil {
		return n, err
	}
	return w.Write([]byte(b))
}

func ReadLd(r io.Reader) (dat []byte, err error) {
	buf := make([]byte, 4)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	tl := binary.BigEndian.Uint32(buf)
	dat = make([]byte, int(tl))
	_, err = io.ReadFull(r, dat)
	return
}
