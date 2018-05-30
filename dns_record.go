package routerosclient

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type resourceDNSStaticRecord struct {
	ID       string `ros:".id"`
	Address  string `ros:"address"  valid:"ipv4,required"`
	Comment  string `ros:"comment"  valid:"optional"`
	Disabled bool   `ros:"disabled" valid:"optional"`
	Name     string `ros:"name"     valid:"required"`
	TTL      string `ros:"ttl"      valid:"optional"`
}

func (d *resourceDNSStaticRecord) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func NewDNSStaticRecord(attrs map[string]string) (*resourceDNSStaticRecord, error) {
	// FIXME
	var d *resourceDNSStaticRecord
	var ok bool

	i, err := setFieldsFromMap(&resourceDNSStaticRecord{}, attrs)
	if err != nil {
		return nil, err
	}

	if d, ok = i.(*resourceDNSStaticRecord); !ok {
		return nil, fmt.Errorf("unable to cast interface to *resourceDNSStaticRecord")
	}

	if err := d.validate(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *resourceDNSStaticRecord) getID() string {
	return d.ID
}

func (d *resourceDNSStaticRecord) setID(id string) {
	d.ID = id
}

func (*resourceDNSStaticRecord) getCreateCommand() string {
	return "/ip/dns/static/add"
}

func (*resourceDNSStaticRecord) getReadCommand() string {
	return "/ip/dns/static/print"
}

func (*resourceDNSStaticRecord) getUpdateCommand() string {
	return "/ip/dns/static/set"
}

func (*resourceDNSStaticRecord) getDeleteCommand() string {
	return "/ip/dns/static/remove"
}
