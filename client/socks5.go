package main

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
)

func handleSocks5(conn net.Conn) (target string, err error) {
	methodCount, err := ReadByte(conn)
	if err != nil {
		return
	}
	methods := make([]byte, methodCount)
	_, err = io.ReadFull(conn, methods)
	if err != nil {
		return
	}
	var hasNAQ bool = false
	for _, n := range methods {
		hasNAQ = n == 0x00
		if hasNAQ {
			break
		}
	}
	if hasNAQ {
		//SELECTED NO AUTHENTICATION REQUIRED
		conn.Write([]byte{socksVer5, NO_AUTHENTICATION_REQUIRED})
	} else {
		conn.Write([]byte{socksVer5, NO_ACCEPTABLE_METHODS})
		return "", errors.New("client not support NAQ")
	}
	ver, err := ReadByte(conn)
	if err != nil {
		return
	}
	if ver != socksVer5 {
		return "", errors.New("socks ver must to be 0x05")
	}
	cmd, err := ReadByte(conn)
	if err != nil {
		return
	}
	if cmd != socksCmdConnect {
		conn.Write([]byte{socksVer5, cmdNotSupport})
		return "", errors.New("not supported command")
	}
	_, err = ReadByte(conn) //skip RSV byte
	if err != nil {
		return
	}
	atyp, err := ReadByte(conn)
	if err != nil {
		return
	}
	var host string
	var port uint16
	if atyp == ATYP_IP4 {
		buf := make([]byte, net.IPv4len+2)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			return
		}
		host = net.IP(buf[:net.IPv4len]).String()
		port = binary.BigEndian.Uint16(buf[net.IPv4len:])
	} else if atyp == ATYP_IP6 {
		buf := make([]byte, net.IPv6len+2)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			return
		}
		host = net.IP(buf[:net.IPv6len]).String()
		port = binary.BigEndian.Uint16(buf[net.IPv6len:])
	} else if atyp == ATYP_DOMAIN {
		domainLength, err := ReadByte(conn)
		if err != nil {
			return "", err
		}
		buf := make([]byte, domainLength+2)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			return "", err
		}
		host = string(buf[0:domainLength])
		port = binary.BigEndian.Uint16(buf[domainLength:])
	} else {
		return "", errors.New("not supported address type")
	}
	target = net.JoinHostPort(host, strconv.Itoa(int(port)))
	if err == nil {
		_, err = conn.Write([]byte{socksVer5, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00})
	}
	return
}
