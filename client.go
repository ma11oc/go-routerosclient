package routerosclient

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/go-routeros/routeros"
)

type TLSConfig struct {
	conf *tls.Config
}

type Config struct {
	Address   string     `valid:"-"`
	Username  string     `valid:"required"`
	Password  string     `valid:"required"`
	Async     bool       `valid:"optional"`
	TLSConfig *TLSConfig `valid:"optional"`
}

type Conn interface {
	RunArgs([]string) (*routeros.Reply, error)
	Close()
}

type Client struct {
	mu   sync.Mutex
	conn Conn
}

func NewClient(c *Config) (*Client, error) {

	if err := c.validate(); err != nil {
		return nil, err
	}

	if c.TLSConfig != nil {
		return NewTLSClient(c)
	} else {
		return NewInsecureClient(c)
	}
}

func NewInsecureClient(c *Config) (*Client, error) {
	rc, err := routeros.Dial(c.Address, c.Username, c.Password)

	if err != nil {
		return nil, err
	}

	return &Client{
		conn: rc,
	}, nil
}

func NewTLSClient(c *Config) (*Client, error) {
	// TODO: validate tls config
	// rc, err := routeros.DialTLS(c.address, c.username, c.password, c.tlsConfig.conf)
	rc, err := routeros.DialTLS(c.Address, c.Username, c.Password, nil)

	if err != nil {
		return nil, err
	}

	return &Client{
		conn: rc,
	}, nil
}

func (c *Client) Run(query string) (*routeros.Reply, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.RunArgs(strings.Split(query, " "))
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Config) validate() error {

	if _, err := govalidator.ValidateStruct(c); err != nil {
		return err
	}

	host, port, err := net.SplitHostPort(c.Address)

	if err != nil {
		return fmt.Errorf("unable to parse RouterOS address")
	}

	if !govalidator.IsIPv4(host) {
		return fmt.Errorf("invalid IPv4 address")
	}

	if !govalidator.IsPort(port) {
		return fmt.Errorf("invalid port")
	}

	if c.Username == "" || c.Password == "" {
		return fmt.Errorf("username and password are required")
	}

	return nil
}
