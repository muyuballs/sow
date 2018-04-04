package main

type Config struct {
	UDT        bool
	Listen     string
	Zlib       bool
	Key        string
	SMux       bool
	SockBuf    int
	Rwnd       int
	Swnd       int
	Mtu        int
	NoDelay    bool
	Server     string
	HttpEnable bool
}
