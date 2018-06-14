package routerosclient

import (
	"fmt"
	"log"
)

type Resource interface {
	validate() error
	getID() string
	setID(string)
	getCreateCommand() string
	getReadCommand() string
	getUpdateCommand() string
	getDeleteCommand() string
}

// CreateResource FIXME: add description
func (c *Client) CreateResource(res Resource) (string, error) {
	log.Printf("[D][C] CreateResource(%v)", res)

	if err := res.validate(); err != nil {
		return "", err
	}

	if err, ok := c.CheckResourceExists(res); ok {
		return "", fmt.Errorf("resource exists: %v", res)
	} else if err != nil {
		return "", err
	}

	command := res.getCreateCommand()
	attrs, err := buildAttrsFromResource(res)
	if err != nil {
		return "", err
	}

	cmd, err := buildCommand(command, nil, &attrs, false)
	if err != nil {
		return "", err
	}
	log.Printf("[D][C][->] %v", cmd)

	r, err := c.Run(cmd)
	if err != nil {
		log.Printf("[E][C][<-] error: %v", err)
		return "", err
	}
	log.Printf("[D][C][<-] %v | %v", r.Re, r.Done)

	if r.Done != nil && r.Done.Map["ret"] != "" {
		return r.Done.Map["ret"], nil
	}

	return "", fmt.Errorf("unexpected empty reply from RouterOS")
}

// UpdateResource ... FIXME
func (c *Client) UpdateResource(o Resource, n Resource) (error, bool) {
	log.Printf("[D][U] UpdateResource")
	log.Printf("[D][U] o(%v) -> n(%v)", o, n)

	if err := o.validate(); err != nil {
		return err, false
	}

	if err := n.validate(); err != nil {
		return err, false
	}

	cur, err := c.ReadResource(o)
	if err != nil {
		return err, false
	}

	n.setID(cur.getID())

	command := o.getUpdateCommand()
	attrs, err := buildAttrsFromResource(n)
	if err != nil {
		return err, false
	}

	cmd, err := buildCommand(command, nil, &attrs, false)
	log.Printf("[D][U][->] %v", cmd)

	r, err := c.Run(cmd)

	if err != nil {
		log.Printf("[E][U][<-] error: %v", err)
		return err, false
	}
	log.Printf("[D][U][<-] %v | %v", r.Re, r.Done)

	return nil, true
}

func (c *Client) ReadResource(res Resource) (Resource, error) {
	log.Printf("[D][R] ReadResource(%v)", res)

	var attrs map[string]string
	var err error

	if err = res.validate(); err != nil {
		return nil, err
	}

	command := res.getReadCommand()

	// it's enough to have a non empty id field to perform query
	if res.getID() != "" {
		attrs = map[string]string{".id": res.getID()}
	} else {
		attrs, err = buildAttrsFromResource(res)
		if err != nil {
			return nil, err
		}
	}

	cmd, err := buildCommand(command, nil, &attrs, true)
	if err != nil {
		return nil, err
	}
	log.Printf("[D][R][->] %v", cmd)

	r, err := c.Run(cmd)
	if err != nil {
		log.Printf("[E][R][<-] error: %v", err)
		return nil, err
	}
	log.Printf("[D][R][<-] %v | %v", r.Re, r.Done)

	switch rlen := len(r.Re); rlen {
	case 0:
		return nil, fmt.Errorf("no resource has been found: %v", res)
	case 1:
		var nr Resource

		switch res.(type) {
		case *ResourceInterfaceBridge:
			nr = &ResourceInterfaceBridge{}
		case *ResourceDHCPServer:
			nr = &ResourceDHCPServer{}
		case *ResourceDHCPServerNetwork:
			nr = &ResourceDHCPServerNetwork{}
		case *ResourceDHCPServerOption:
			nr = &ResourceDHCPServerOption{}
		case *ResourceDHCPServerOptionSet:
			nr = &ResourceDHCPServerOptionSet{}
		case *ResourceDHCPServerLease:
			nr = &ResourceDHCPServerLease{}
		case *ResourceDNSStaticRecord:
			nr = &ResourceDNSStaticRecord{}
		default:
			return nil, fmt.Errorf("unable to determine resource type")
		}

		obj, err := setFieldsFromMap(nr, r.Re[0].Map)
		if err != nil {
			return nil, err
		}

		return obj, nil

	default:
		return nil, fmt.Errorf("ambiguous reply")
	}
}

// DeleteResource deletes existing lease from RouterOS. If lease doesn't exist, returns error.
func (c *Client) DeleteResource(res Resource) (error, bool) {
	log.Printf("[D][D] DeleteResource(%v)", res)

	if err := res.validate(); err != nil {
		return err, false
	}

	resource, err := c.ReadResource(res)
	if err != nil {
		return err, false
	}

	command := res.getDeleteCommand()
	attrs := map[string]string{".id": resource.getID()}

	cmd, err := buildCommand(command, nil, &attrs, false)
	if err != nil {
		return err, false
	}
	log.Printf("[D][D][->] %v", cmd)

	r, err := c.Run(cmd)
	if err != nil {
		log.Printf("[E][D][<-] error: %v", err)
		return err, false
	}
	log.Printf("[D][D][<-] %v | %v", r.Re, r.Done)

	return nil, true
}

func (c *Client) CheckResourceExists(res Resource) (error, bool) {
	log.Printf("[D][?] CheckResourceExists(%v)", res)

	if err := res.validate(); err != nil {
		return err, false
	}

	command := res.getReadCommand()
	proplist := []string{".id"}

	attrs, err := buildAttrsFromResource(res)
	if err != nil {
		return err, false
	}

	cmd, err := buildCommand(command, &proplist, &attrs, true)
	if err != nil {
		return err, false
	}
	log.Printf("[D][?][->] %v", cmd)

	r, err := c.Run(cmd)
	if err != nil {
		log.Printf("[E][?][<-] error: %v", err)
		return err, false
	}
	log.Printf("[D][?][<-] %v | %v", r.Re, r.Done)

	switch rlen := len(r.Re); rlen {
	case 0:
		return nil, false
	case 1:
		return nil, true
	default:
		return fmt.Errorf("ambiguous reply: %v", r), false
	}
}
