package core

import (
	"io"
)

type Config struct {
	UDT        bool
	Listen     string
	Key        string
	SMux       bool
	SockBuf    int
	Rwnd       int
	Swnd       int
	Mtu        int
	NoDelay    bool
	Server     string
	HttpEnable bool
	Momo       bool
	MomoAddr   string
	LogFile    string
	LogOut     io.Writer
	LOG_FLAGS  int
}
