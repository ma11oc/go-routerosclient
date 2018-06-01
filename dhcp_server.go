package routerosclient

import (
	"github.com/asaskevich/govalidator"
)

type ResourceDHCPServer struct {
	ID        string `ros:".id"`
	Disabled  bool   `ros:"disabled"     valid:"optional"`
	Name      string `ros:"name"         valid:"optional"`
	Interface string `ros:"interface"    valid:"required"`
}

func (d *ResourceDHCPServer) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func (d *ResourceDHCPServer) getID() string {
	return d.ID
}

func (d *ResourceDHCPServer) setID(id string) {
	d.ID = id
}

func (*ResourceDHCPServer) getCreateCommand() string {
	return "/ip/dhcp-server/add"
}

func (*ResourceDHCPServer) getReadCommand() string {
	return "/ip/dhcp-server/print"
}

func (*ResourceDHCPServer) getUpdateCommand() string {
	return "/ip/dhcp-server/set"
}

func (*ResourceDHCPServer) getDeleteCommand() string {
	return "/ip/dhcp-server/remove"
}
