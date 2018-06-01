package routerosclient

import (
	"github.com/asaskevich/govalidator"
)

/*
 * TODO:
 * admin-mac --
 * ageing-time -- Time the information about host will be kept in the the data base
 * arp -- Address Resolution Protocol
 * arp-timeout --
 * auto-mac --
 * max-message-age -- Time to remember Hello messages received from other bridges
 * priority -- Bridge interface priority
 * protocol-mode --
 * transmit-hold-count --
 */

type ResourceInterfaceBridge struct {
	ID           string `ros:".id"`
	Comment      string `ros:"comment"       valid:"optional"`
	Disabled     bool   `ros:"disabled"      valid:"optional"`
	FastForward  bool   `ros:"fast-forward"  valid:"optional"`
	ForwardDelay string `ros:"forward-delay" valid:"optional"`
	MTU          int    `ros:"mtu"           valid:"optional"`
	Name         string `ros:"name"          valid:"required"`
}

func (d *ResourceInterfaceBridge) validate() error {
	if d.ID == "" {
		_, err := govalidator.ValidateStruct(d)

		if err != nil {
			return err
		}

	}

	return nil
}

func (d *ResourceInterfaceBridge) getID() string {
	return d.ID
}

func (d *ResourceInterfaceBridge) setID(id string) {
	d.ID = id
}

func (*ResourceInterfaceBridge) getCreateCommand() string {
	return "/interface/bridge/add"
}

func (*ResourceInterfaceBridge) getReadCommand() string {
	return "/interface/bridge/print"
}

func (*ResourceInterfaceBridge) getUpdateCommand() string {
	return "/interface/bridge/set"
}

func (*ResourceInterfaceBridge) getDeleteCommand() string {
	return "/interface/bridge/remove"
}
