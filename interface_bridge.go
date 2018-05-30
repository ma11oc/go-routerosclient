package routerosclient

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

/*
 * TODO:
 * admin-mac --
 * ageing-time -- Time the information about host will be kept in the the data base
 * arp -- Address Resolution Protocol
 * arp-timeout --
 * auto-mac --
 * comment -- Short description of the item
 * disabled -- Defines whether item is ignored or used
 * fast-forward --
 * forward-delay -- Time which is spent in listening/learning state
 * max-message-age -- Time to remember Hello messages received from other bridges
 * mtu -- Maximum Transmit Unit
 * name -- Bridge name
 * priority -- Bridge interface priority
 * protocol-mode --
 * transmit-hold-count --
 */

type resourceInterfaceBridge struct {
	ID           string `ros:".id"`
	Comment      string `ros:"comment"       valid:"optional"`
	Disabled     bool   `ros:"disabled"      valid:"optional"`
	FastForward  bool   `ros:"fast-forward"  valid:"optional"`
	ForwardDelay string `ros:"forward-delay" valid:"optional"`
	MTU          int    `ros:"mtu"           valid:"optional"`
	Name         string `ros:"name"          valid:"required"`
}

func (d *resourceInterfaceBridge) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func NewInterfaceBridge(attrs map[string]string) (*resourceInterfaceBridge, error) {
	// FIXME
	var d *resourceInterfaceBridge
	var ok bool

	i, err := setFieldsFromMap(&resourceInterfaceBridge{}, attrs)
	if err != nil {
		return nil, err
	}

	if d, ok = i.(*resourceInterfaceBridge); !ok {
		return nil, fmt.Errorf("unable to cast interface to *resourceInterfaceBridge")
	}

	if err := d.validate(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *resourceInterfaceBridge) getID() string {
	return d.ID
}

func (d *resourceInterfaceBridge) setID(id string) {
	d.ID = id
}

func (*resourceInterfaceBridge) getCreateCommand() string {
	return "/interface/bridge/add"
}

func (*resourceInterfaceBridge) getReadCommand() string {
	return "/interface/bridge/print"
}

func (*resourceInterfaceBridge) getUpdateCommand() string {
	return "/interface/bridge/set"
}

func (*resourceInterfaceBridge) getDeleteCommand() string {
	return "/interface/bridge/remove"
}
