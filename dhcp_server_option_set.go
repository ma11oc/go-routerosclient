package routerosclient

import (
	"github.com/asaskevich/govalidator"
)

type ResourceDHCPServerOptionSet struct {
	ID      string `ros:".id"`
	Name    string `ros:"name"     valid:"required"`
	Options string `ros:"options"  valid:"required"`
}

func (d *ResourceDHCPServerOptionSet) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func (d *ResourceDHCPServerOptionSet) getID() string {
	return d.ID
}

func (d *ResourceDHCPServerOptionSet) setID(id string) {
	d.ID = id
}

func (*ResourceDHCPServerOptionSet) getCreateCommand() string {
	return "/ip/dhcp-server/option/sets/add"
}

func (*ResourceDHCPServerOptionSet) getReadCommand() string {
	return "/ip/dhcp-server/option/sets/print"
}

func (*ResourceDHCPServerOptionSet) getUpdateCommand() string {
	return "/ip/dhcp-server/option/sets/set"
}

func (*ResourceDHCPServerOptionSet) getDeleteCommand() string {
	return "/ip/dhcp-server/option/sets/remove"
}
