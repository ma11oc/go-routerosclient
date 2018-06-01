package routerosclient

import (
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
type ResourceDHCPServerLease struct {
	ID           string `ros:".id"`
	Address      string `ros:"address"           valid:"ipv4,required"`
	AddressLists string `ros:"address-lists"     valid:"optional"`
	ClientID     string `ros:"client-id"         valid:"optional"`
	Comment      string `ros:"comment"           valid:"optional"`
	Disabled     bool   `ros:"disabled"          valid:"optional"`
	MacAddress   string `ros:"mac-address"       valid:"mac,required"`
	Server       string `ros:"server"            valid:"required"`
}

func (d *ResourceDHCPServerLease) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

/* FIXME
 * func NewDHCPServerLease(attrs map[string]string) (*ResourceDHCPServerLease, error) {
 *     var d *ResourceDHCPServerLease
 *     var ok bool
 *
 *     i, err := setFieldsFromMap(&resourceDHCPServerLease{}, attrs)
 *     if err != nil {
 *         return nil, err
 *     }
 *
 *     if d, ok = i.(*ResourceDHCPServerLease); !ok {
 *         return nil, fmt.Errorf("unable to cast interface to *ResourceDHCPServerLease")
 *     }
 *
 *     if err := d.validate(); err != nil {
 *         return nil, err
 *     }
 *
 *     return d, nil
 * }
 */

func (d *ResourceDHCPServerLease) getID() string {
	return d.ID
}

func (d *ResourceDHCPServerLease) setID(id string) {
	d.ID = id
}

func (*ResourceDHCPServerLease) getCreateCommand() string {
	return "/ip/dhcp-server/lease/add"
}

func (*ResourceDHCPServerLease) getReadCommand() string {
	return "/ip/dhcp-server/lease/print"
}

func (*ResourceDHCPServerLease) getUpdateCommand() string {
	return "/ip/dhcp-server/lease/set"
}

func (*ResourceDHCPServerLease) getDeleteCommand() string {
	return "/ip/dhcp-server/lease/remove"
}
