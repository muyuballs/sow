package main

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
)

func skipIDEN(conn net.Conn) error {
	for {
		b, err := ReadByte(conn)
		if err != nil {
			return err
		}
		if b == 0x00 {
			return nil
		}
	}
}

func handleSocks4(conn net.Conn) (target string, err error) {
	cmd, err := ReadByte(conn)
	if err != nil {
		return
	}
	if cmd != socksCmdConnect {
		conn.Write([]byte{socksVer4, REQUEST_REJECTED})
		return "", errors.New("not supported command")
	}
	buf := make([]byte, 2)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return
	}
	port := binary.BigEndian.Uint16(buf)
	buf = make([]byte, net.IPv4len)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return
	}
	if buf[0] == 0x00 && buf[1] == 0x00 && buf[2] == 0x00 { //socks4a ip address
		err = skipIDEN(conn)
		if err != nil {
			return
		}
		buf = make([]byte, 0)
		for {
			b, err := ReadByte(conn)
			if err != nil {
				return "", err
			}
			if b == 0x00 {
				break
			}
			buf = append(buf, b)
		}
		target = net.JoinHostPort(string(buf), strconv.Itoa(int(port)))
	} else {
		target = net.JoinHostPort(net.IP(buf).String(), strconv.Itoa(int(port)))
		err = skipIDEN(conn)
	}
	if err == nil {
		_, err = conn.Write([]byte{socksVer4, 0x5A})
	}
	return
}
