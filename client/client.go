package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/urfave/cli"
)

var (
	VERSION = "0.0.2"
)

var transferFn func(string, *net.TCPConn, *Config)

func handleConnection(conn *net.TCPConn, c *Config) {
	err := conn.SetNoDelay(true)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	ver, err := ReadByte(conn)
	if err != nil {
		return
	}
	log.Println("client version ", ver)
	if ver == socksVer5 {
		target, err := handleSocks5(conn)
		if err != nil {
			log.Println(err)
			return
		}
		transferFn(target, conn, c)
	} else if ver == socksVer4 {
		target, err := handleSocks4(conn)
		if err != nil {
			log.Println(err)
			return
		}
		transferFn(target, conn, c)
	} else {
		conn.Close()
		log.Println("version", ver, "not supported")
	}
}

func main() {
	myApp := cli.NewApp()
	myApp.Name = "SOW"
	myApp.Usage = "client"
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Value: ":1221",
			Usage: "sow client socks(4,5) listen address",
		},
		cli.StringFlag{
			Name:  "server, s",
			Value: "",
			Usage: "sow server address",
		},
		cli.StringFlag{
			Name:  "key, k",
			Value: "qwert",
			Usage: "secret key",
		},
		cli.BoolFlag{
			Name:  "zlib, z",
			Usage: "use zlib compress data",
		},
		cli.BoolFlag{
			Name:  "udt, u",
			Usage: "use udp ",
		},
		cli.BoolFlag{
			Name:  "smux",
			Usage: "tcp with smux",
		},
		cli.IntFlag{
			Name:   "sockBuf",
			Value:  4194304,
			Hidden: true,
		},
		cli.IntFlag{
			Name:   "rwnd",
			Value:  1024,
			Hidden: true,
		},
		cli.IntFlag{
			Name:   "swnd",
			Value:  1024,
			Hidden: true,
		}, cli.IntFlag{
			Name:   "mtu",
			Value:  1350,
			Hidden: true,
		},
		cli.BoolFlag{
			Name:  "nodelay, nd",
			Usage: "tcp nodelay, udt ack nodelay",
		},
	}

	myApp.Action = func(c *cli.Context) error {
		config := &Config{}
		config.Server = c.String("server")
		config.Key = c.String("key")
		config.Listen = c.String("listen")
		config.Zlib = c.Bool("zlib")
		config.UDT = c.Bool("udt")
		config.SMux = c.Bool("smux")
		config.SockBuf = c.Int("sockBuf")
		config.Rwnd = c.Int("rwnd")
		config.Swnd = c.Int("swnd")
		config.Mtu = c.Int("mtu")
		config.NoDelay = c.Bool("nodelay")
		if config.Server == "" {
			return errors.New("server address is null")
		}
		localAddr, err := net.ResolveTCPAddr("tcp", config.Listen)
		if err != nil {
			return err
		}
		ln, err := net.ListenTCP("tcp", localAddr)
		if err != nil {
			return err
		}
		defer ln.Close()
		if config.UDT {
			transferFn = transferByUDT
		} else {
			transferFn = transferByTCP
		}
		log.Println("Start @", ln.Addr())
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				return err
			}
			log.Println("new client:", conn.RemoteAddr())
			go handleConnection(conn, config)
		}
	}
	log.Println(myApp.Run(os.Args))
}
