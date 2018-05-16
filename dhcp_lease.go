package routerosclient

import (
	"fmt"
	"log"

	"github.com/asaskevich/govalidator"
)

// DHCPLease is a struct which describes dhcp lease.
// `ros` tag contains valid attributes names from the RouterOS point of view.
// TODO: insert-queue-before
// TODO: lease-time
// TODO: rate-limit
// BUG: Surprisingly, RouterOS expects `=blocked=bool` on writing and `?=block-access=bool` on reading for `block-access` attribute.
// BUG: RouterOS does not recognize space separated value of `comment` attribute.
// BUG: AlwaysBroadcast doesn't use for read query.
// BUG: UseSrcMac doesn't use for read query.
type DHCPLease struct {
	ID            string `ros:".id"`
	Address       string `ros:"address"           valid:"ipv4,required"`
	AddressLists  string `ros:"address-lists"     valid:"optional"`
	ClientID      string `ros:"client-id"         valid:"optional"`
	Comment       string `ros:"comment"           valid:"optional"`
	DhcpOption    string `ros:"dhcp-option"       valid:"optional"`
	DhcpOptionSet string `ros:"dhcp-option-set"   valid:"optional"`
	Disabled      bool   `ros:"disabled"          valid:"optional"`
	MacAddress    string `ros:"mac-address"       valid:"mac,required"`
	Server        string `ros:"server"            valid:"required"`
}

func (d *DHCPLease) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

// CreateDHCPLease ...
// FIXME: add description
func (c *Client) CreateDHCPLease(d *DHCPLease) (string, error) {
	log.Printf("[D][C] CreateDHCPLease(%v)", d)

	if err := d.validate(); err != nil {
		return "", err
	}

	if err, ok := c.CheckDHCPLeaseExists(d); ok {
		return "", fmt.Errorf("dhcp lease exists")
	} else if err != nil {
		return "", err
	}

	command := "/ip/dhcp-server/lease/add"

	attrs, err := buildAttrsFromResource(d)
	if err != nil {
		return "", err
	}

	cmd, err := buildCommand(command, nil, &attrs, false)
	if err != nil {
		return "", err
	}
	log.Printf("[D][C][->] %v", cmd)

	r, err := c.Run(cmd)
	if err != nil {
		return "", err
	}
	log.Printf("[D][C][<-] %v", r)

	if r.Done != nil && r.Done.Map["ret"] != "" {
		return r.Done.Map["ret"], nil
	}

	return "", fmt.Errorf("unexpected empty reply from RouterOS")
}

// UpdateDHCPLease ... FIXME
func (c *Client) UpdateDHCPLease(o *DHCPLease, n *DHCPLease) (error, bool) {
	log.Printf("[D][U] UpdateDHCPLease")
	log.Printf("[D][U] o(%v) -> n(%v)", o, n)

	if err := o.validate(); err != nil {
		return err, false
	}

	if err := n.validate(); err != nil {
		return err, false
	}

	command := "/ip/dhcp-server/lease/set"

	cur, err := c.ReadDHCPLease(o)
	if err != nil {
		return err, false
	}

	n.ID = cur.ID

	attrs, err := buildAttrsFromResource(n)
	if err != nil {
		return err, false
	}

	cmd, err := buildCommand(command, nil, &attrs, false)
	log.Printf("[D][U][->] %v", cmd)

	r, err := c.Run(cmd)

	if err != nil {
		return err, false
	}
	log.Printf("[D][U][<-] %v", r)

	return nil, true
}

func (c *Client) ReadDHCPLease(d *DHCPLease) (*DHCPLease, error) {
	log.Printf("[D][R] ReadDHCPLease(%v)", d)

	if err := d.validate(); err != nil {
		return nil, err
	}

	command := "/ip/dhcp-server/lease/print"

	attrs, err := buildAttrsFromResource(d)
	if err != nil {
		return nil, err
	}

	cmd, err := buildCommand(command, nil, &attrs, true)
	if err != nil {
		return nil, err
	}
	log.Printf("[D][R][->] %v", cmd)

	r, err := c.Run(cmd)
	if err != nil {
		return nil, err
	}
	log.Printf("[D][R][<-] %v", r)

	switch rlen := len(r.Re); rlen {
	case 0:
		return nil, fmt.Errorf("no lease has been found")
	case 1:
		obj, err := setFieldsFromReply(&DHCPLease{}, r.Re[0])
		if err != nil {
			return nil, err
		}

		if newLease, ok := obj.(*DHCPLease); ok {
			return newLease, nil
		}

		return nil, fmt.Errorf("unable to cast interface to *DHCPLease")
	default:
		return nil, fmt.Errorf("ambiguous reply")
	}
}

// DeleteDHCPLease deletes existing lease from RouterOS. Is lease doesn't exist, returns error.
func (c *Client) DeleteDHCPLease(d *DHCPLease) (error, bool) {
	log.Printf("[D][D] DeleteDHCPLease(%v)", d)

	var l *DHCPLease

	if err := d.validate(); err != nil {
		return err, false
	}

	l, err := c.ReadDHCPLease(d)
	if err != nil {
		return err, false
	}

	command := "/ip/dhcp-server/lease/remove"

	attrs := map[string]string{".id": l.ID}

	cmd, err := buildCommand(command, nil, &attrs, false)
	if err != nil {
		return err, false
	}
	log.Printf("[D][D][->] %v", cmd)

	r, err := c.Run(cmd)
	if err != nil {
		return err, false
	}
	log.Printf("[D][D][<-] %v", r)

	return nil, true
}

func (c *Client) CheckDHCPLeaseExists(d *DHCPLease) (error, bool) {
	log.Printf("[D][?] CheckDHCPLeaseExists(%v)", d)

	if err := d.validate(); err != nil {
		return err, false
	}

	command := "/ip/dhcp-server/lease/print"
	proplist := []string{".id"}

	attrs, err := buildAttrsFromResource(d)
	if err != nil {
		return err, false
	}

	cmd, err := buildCommand(command, &proplist, &attrs, true)
	if err != nil {
		return err, false
	}
	log.Printf("[D][?][->] %v", cmd)

	r, err := c.Run(cmd)
	if err != nil {
		return err, false
	}
	log.Printf("[D][?][<-] %v", r)

	switch rlen := len(r.Re); rlen {
	case 0:
		return nil, false
	case 1:
		return nil, true
	default:
		return fmt.Errorf("ambiguous reply: %v", r), false
	}
}
