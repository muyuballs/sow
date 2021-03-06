package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/muyuballs/sow/core"
	"github.com/urfave/cli"
)

var (
	VERSION = "0.0.3"
)

var transferFn func(string, *net.TCPConn, *core.Config)

func handleConnection(conn *net.TCPConn, c *core.Config) {
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
		if c.HttpEnable {
			handleHttp(ver, conn, c)
		} else {
			conn.Close()
			log.Println("not supported ver", ver)
		}
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
			Usage: "sow client socks(4,5) & http listen address",
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
			Name:  "udt, u",
			Usage: "use udp ",
		},
		cli.StringFlag{
			Name:  "log",
			Usage: "log file",
			Value: "console",
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
		cli.BoolFlag{
			Name:  "disable-http, dh",
			Usage: "disable http proxy",
		},
		cli.BoolFlag{
			Name:  "momo",
			Usage: "enable momo server",
		},
		cli.StringFlag{
			Name:  "momo-addr",
			Value: ":60001",
			Usage: "momo server listen addr",
		},
	}

	myApp.Action = func(c *cli.Context) error {
		config := &core.Config{}
		config.LOG_FLAGS = log.LstdFlags | log.Lmicroseconds
		config.Server = c.String("server")
		config.Key = c.String("key")
		config.Listen = c.String("listen")
		config.UDT = c.Bool("udt")
		config.SMux = c.Bool("smux")
		config.SockBuf = c.Int("sockBuf")
		config.Rwnd = c.Int("rwnd")
		config.Swnd = c.Int("swnd")
		config.Mtu = c.Int("mtu")
		config.NoDelay = c.Bool("nodelay")
		config.HttpEnable = !c.Bool("disable-http")
		config.Momo = c.Bool("momo")
		config.MomoAddr = c.String("momo-addr")
		config.LogFile = c.String("log")
		if "console" != config.LogFile {
			logOut, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0755)
			if err != nil {
				return err
			}
			defer logOut.Close()
			defer logOut.Sync()
			config.LogOut = logOut
		} else {
			config.LogOut = os.Stdout
		}
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
		if config.Momo {
			mlAddr, err := net.ResolveTCPAddr("tcp", config.MomoAddr)
			if err != nil {
				return err
			}
			mln, err := net.ListenTCP("tcp", mlAddr)
			if err != nil {
				return err
			}
			log.Println("Start Momo Server @:", mln.Addr())
			go func() {
				for {
					c, err := mln.AcceptTCP()
					if err != nil {
						log.Println(err)
						return
					}
					go handMomoConnection(c, config)
				}
			}()
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
