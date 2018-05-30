package routerosclient

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

// resourceDHCPServerLease is a struct which describes dhcp lease.
// `ros` tag contains valid attributes names from the RouterOS point of view.
// TODO: insert-queue-before
// TODO: lease-time
// TODO: rate-limit
// TODO: DHCPOption    string `ros:"dhcp-option"       valid:"optional"`
// TODO: DHCPOptionSet string `ros:"dhcp-option-set"   valid:"optional"`
// BUG: Surprisingly, RouterOS expects `=blocked=bool` on writing and `?=block-access=bool` on reading for `block-access` attribute.
// BUG: RouterOS does not recognize space separated value of `comment` attribute.
// BUG: AlwaysBroadcast doesn't use for read query.
// BUG: UseSrcMac doesn't use for read query.
type resourceDHCPServerLease struct {
	ID           string `ros:".id"`
	Address      string `ros:"address"           valid:"ipv4,required"`
	AddressLists string `ros:"address-lists"     valid:"optional"`
	ClientID     string `ros:"client-id"         valid:"optional"`
	Comment      string `ros:"comment"           valid:"optional"`
	Disabled     bool   `ros:"disabled"          valid:"optional"`
	MacAddress   string `ros:"mac-address"       valid:"mac,required"`
	Server       string `ros:"server"            valid:"required"`
}

func (d *resourceDHCPServerLease) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func NewDHCPServerLease(attrs map[string]string) (*resourceDHCPServerLease, error) {
	// FIXME
	var d *resourceDHCPServerLease
	var ok bool

	i, err := setFieldsFromMap(&resourceDHCPServerLease{}, attrs)
	if err != nil {
		return nil, err
	}

	if d, ok = i.(*resourceDHCPServerLease); !ok {
		return nil, fmt.Errorf("unable to cast interface to *resourceDHCPServerLease")
	}

	if err := d.validate(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *resourceDHCPServerLease) getID() string {
	return d.ID
}

func (d *resourceDHCPServerLease) setID(id string) {
	d.ID = id
}

func (*resourceDHCPServerLease) getCreateCommand() string {
	return "/ip/dhcp-server/lease/add"
}

func (*resourceDHCPServerLease) getReadCommand() string {
	return "/ip/dhcp-server/lease/print"
}

func (*resourceDHCPServerLease) getUpdateCommand() string {
	return "/ip/dhcp-server/lease/set"
}

func (*resourceDHCPServerLease) getDeleteCommand() string {
	return "/ip/dhcp-server/lease/remove"
}
