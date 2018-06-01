package routerosclient

import (
	"github.com/asaskevich/govalidator"
)

type ResourceDNSStaticRecord struct {
	ID       string `ros:".id"`
	Address  string `ros:"address"  valid:"ipv4,required"`
	Comment  string `ros:"comment"  valid:"optional"`
	Disabled bool   `ros:"disabled" valid:"optional"`
	Name     string `ros:"name"     valid:"required"`
	TTL      string `ros:"ttl"      valid:"optional"`
}

func (d *ResourceDNSStaticRecord) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func (d *ResourceDNSStaticRecord) getID() string {
	return d.ID
}

func (d *ResourceDNSStaticRecord) setID(id string) {
	d.ID = id
}

func (*ResourceDNSStaticRecord) getCreateCommand() string {
	return "/ip/dns/static/add"
}

func (*ResourceDNSStaticRecord) getReadCommand() string {
	return "/ip/dns/static/print"
}

func (*ResourceDNSStaticRecord) getUpdateCommand() string {
	return "/ip/dns/static/set"
}

func (*ResourceDNSStaticRecord) getDeleteCommand() string {
	return "/ip/dns/static/remove"
}
