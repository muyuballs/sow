package main

import (
	"os"

	"github.com/urfave/cli"
)

var (
	VERSION = "0.0.2"
)

func main() {
	myApp := cli.NewApp()
	myApp.Name = "SOW"
	myApp.Usage = "server"
	myApp.Version = VERSION
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Value: ":60000",
			Usage: "sow server listen address",
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
			Name:  "smux, s",
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
		if config.UDT {
			return handleUDT(config)
		}
		return handleTCP(config)
	}
	myApp.Run(os.Args)
}
