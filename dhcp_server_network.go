package routerosclient

import (
	"github.com/asaskevich/govalidator"
)

type ResourceDHCPServerNetwork struct {
	ID            string `ros:".id"`
	Address       string `ros:"address"         valid:"optional"`
	BootFileName  string `ros:"boot-file-name"  valid:"optional"`
	Comment       string `ros:"comment"         valid:"optional"`
	DHCPOption    string `ros:"dhcp-option"     valid:"optional"`
	DHCPOptionSet string `ros:"dhcp-option-set" valid:"optional"`
	Domain        string `ros:"domain"          valid:"dns,optional"`
	DNSServer     string `ros:"dns-server"      valid:"ipv4,optional"`
	Gateway       string `ros:"gateway"         valid:"ipv4,optional"`
	Netmask       string `ros:"netmask"         valid:"optional"`
	NextServer    string `ros:"next-server"     valid:"ipv4,optional"`
	NTPServer     string `ros:"ntp-server"      valid:"ipv4,optional"`
	WINSServer    string `ros:"wins-server"     valid:"ipv4,optional"`
}

func (d *ResourceDHCPServerNetwork) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func (d *ResourceDHCPServerNetwork) getID() string {
	return d.ID
}

func (d *ResourceDHCPServerNetwork) setID(id string) {
	d.ID = id
}

func (*ResourceDHCPServerNetwork) getCreateCommand() string {
	return "/ip/dhcp-server/network/add"
}

func (*ResourceDHCPServerNetwork) getReadCommand() string {
	return "/ip/dhcp-server/network/print"
}

func (*ResourceDHCPServerNetwork) getUpdateCommand() string {
	return "/ip/dhcp-server/network/set"
}

func (*ResourceDHCPServerNetwork) getDeleteCommand() string {
	return "/ip/dhcp-server/network/remove"
}
