package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/dmichael/go-multicast/multicast"
	"github.com/urfave/cli"
)

const (
	defaultMulticastAddress = "239.0.0.0:9999"
)

func main() {
	pl := pingoListener{}
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		address := c.Args().Get(0)
		if address == "" {
			address = defaultMulticastAddress
		}
		fmt.Printf("Listening on %s\n", address)
		go pl.ping(defaultMulticastAddress)
		multicast.Listen(address, pl.msgHandler)
		return nil
	}

	app.Run(os.Args)
}

type pingoListener struct {
	localIPAddress string
}

func (pl *pingoListener) msgHandler(src *net.UDPAddr, n int, b []byte) {
	if pl.localIPAddress == src.String() {
		return
	}
	log.Println(n, "bytes read from", src)
	log.Println(hex.Dump(b[:n]))
}

func (pl *pingoListener) ping(addr string) {
	conn, err := multicast.NewBroadcaster(addr)
	pl.localIPAddress = conn.LocalAddr().String()
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn.Write([]byte("hello, world\n"))
		time.Sleep(1 * time.Second)
	}
}
