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

/* FIXME
 * func NewDHCPServer(attrs map[string]string) (*ResourceDHCPServer, error) {
 *     var d *ResourceDHCPServer
 *     var ok bool
 *
 *     i, err := setFieldsFromMap(&resourceDHCPServer{}, attrs)
 *     if err != nil {
 *         return nil, err
 *     }
 *
 *     if d, ok = i.(*ResourceDHCPServer); !ok {
 *         return nil, fmt.Errorf("unable to cast interface to *ResourceDHCPServer")
 *     }
 *
 *     if err := d.validate(); err != nil {
 *         return nil, err
 *     }
 *
 *     return d, nil
 * }
 */

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
