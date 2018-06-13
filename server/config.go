package main

import (
	"io"
)

type Config struct {
	UDT     bool
	Listen  string
	Zlib    bool
	Key     string
	SMux    bool
	SockBuf int
	Rwnd    int
	Swnd    int
	Mtu     int
	NoDelay bool
	LogFile string
	LogOut  io.Writer
}
