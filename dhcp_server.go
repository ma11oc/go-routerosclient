package routerosclient

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type resourceDHCPServer struct {
	ID        string `ros:".id"`
	Disabled  bool   `ros:"disabled"     valid:"optional"`
	Name      string `ros:"name"         valid:"optional"`
	Interface string `ros:"interface"    valid:"required"`
}

func (d *resourceDHCPServer) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func NewDHCPServer(attrs map[string]string) (*resourceDHCPServer, error) {
	var d *resourceDHCPServer
	var ok bool

	i, err := setFieldsFromMap(&resourceDHCPServer{}, attrs)
	if err != nil {
		return nil, err
	}

	if d, ok = i.(*resourceDHCPServer); !ok {
		return nil, fmt.Errorf("unable to cast interface to *resourceDHCPServer")
	}

	if err := d.validate(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *resourceDHCPServer) getID() string {
	return d.ID
}

func (d *resourceDHCPServer) setID(id string) {
	d.ID = id
}

func (*resourceDHCPServer) getCreateCommand() string {
	return "/ip/dhcp-server/add"
}

func (*resourceDHCPServer) getReadCommand() string {
	return "/ip/dhcp-server/print"
}

func (*resourceDHCPServer) getUpdateCommand() string {
	return "/ip/dhcp-server/set"
}

func (*resourceDHCPServer) getDeleteCommand() string {
	return "/ip/dhcp-server/remove"
}
