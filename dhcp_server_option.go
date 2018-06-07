package routerosclient

import (
	"github.com/asaskevich/govalidator"
)

/*
 * NOTE: string values must be surrounded by quotes:
 * &ResourceDHCPServerOption{
 *     Code: 66,
 *     Name: "next-server",
 *     Value: "'192.168.0.2'"
 * }
 * For details see: https://wiki.mikrotik.com/wiki/Manual:IP/DHCP_Server#DHCP_Options
 *
 */
type ResourceDHCPServerOption struct {
	ID    string `ros:".id"`
	Code  int    `ros:"code"   valid:"required"`
	Name  string `ros:"name"   valid:"required"`
	Value string `ros:"value"  valid:"optional"`
}

func (d *ResourceDHCPServerOption) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func (d *ResourceDHCPServerOption) getID() string {
	return d.ID
}

func (d *ResourceDHCPServerOption) setID(id string) {
	d.ID = id
}

func (*ResourceDHCPServerOption) getCreateCommand() string {
	return "/ip/dhcp-server/option/add"
}

func (*ResourceDHCPServerOption) getReadCommand() string {
	return "/ip/dhcp-server/option/print"
}

func (*ResourceDHCPServerOption) getUpdateCommand() string {
	return "/ip/dhcp-server/option/set"
}

func (*ResourceDHCPServerOption) getDeleteCommand() string {
	return "/ip/dhcp-server/option/remove"
}
