package routerosclient

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-routeros/routeros"
)

var (
	connType  string
	connAddr  string
	connPort  uint64
	connUser  string
	connPass  string
	connAsync bool
	connTLS   bool
	silent    bool
)

type Scenario struct {
	conn *ConnStub
}

func init() {
	flag.StringVar(&connType, "conn", "stub", "type of connection <stub|routeros>")
	flag.StringVar(&connAddr, "conn.addr", "127.0.0.1", "address of connection to RouterOS")
	flag.Uint64Var(&connPort, "conn.port", 8728, "port of connection to RouterOS")
	flag.StringVar(&connUser, "conn.user", "vagrant", "RouterOS username")
	flag.StringVar(&connPass, "conn.pass", "vagrant", "RouterOS password")
	flag.BoolVar(&connAsync, "conn.async", false, "use async code")
	flag.BoolVar(&connTLS, "conn.tls", false, "use TLS encrypted connection (usually port is 8729)")
	flag.BoolVar(&silent, "silent", false, "suppress logging output")

	flag.Parse()

	conntypes := map[string]bool{"stub": true, "routeros": true}

	if _, ok := conntypes[connType]; !ok {
		fmt.Printf("unknown connection type: %v\n", connType)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if silent {
		log.SetOutput(ioutil.Discard)
	}
}

func getTestClient() (*Client, error) {
	var c *Client
	var err error

	if connType == "routeros" {
		c, err = NewClient(&Config{
			Address:  fmt.Sprintf("%v:%v", connAddr, connPort),
			Username: connUser,
			Password: connPass,
			Async:    connAsync,
		})

		if err != nil {
			return nil, err
		}
	} else {
		c = &Client{
			conn: &ConnStub{
				repliesQueue: make(chan *routeros.Reply, 1),
			},
		}
	}

	return c, nil
}
