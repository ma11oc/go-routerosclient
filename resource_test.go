package routerosclient

import (
	"fmt"
	"reflect"
	"testing"
)

type testResource struct {
	env  []Resource // slice of required resources
	min  Resource   // resource minimal config
	full Resource   // resource full config
}

func (tr *testResource) setup(c *Client, stub bool) (error, bool) {
	if !stub {
		for _, r := range tr.env {
			if _, err := c.CreateResource(r); err != nil {
				return err, false
			}
		}
	}

	return nil, true
}

func (tr *testResource) teardown(c *Client, stub bool) (error, bool) {
	if !stub {
		// since resources are dependent from each other,
		// delete them in reverse order
		for i := len(tr.env) - 1; i >= 0; i-- {
			if err, ok := c.DeleteResource(tr.env[i]); !ok {
				return err, false
			}
		}
	}

	return nil, true
}

var (
	resources = []*testResource{
		&testResource{
			min: &ResourceInterfaceBridge{
				Disabled: true,
				Name:     "br0",
			},
			full: &ResourceInterfaceBridge{
				Comment:      "default bridge",
				Disabled:     false,
				FastForward:  true,
				ForwardDelay: "30s",
				MTU:          1500,
				Name:         "br0",
			},
		},
		&testResource{
			env: []Resource{
				&ResourceInterfaceBridge{
					Name:     "test-bridge",
					Disabled: false,
					MTU:      1500,
				},
			},
			min: &ResourceDHCPServer{
				Interface: "test-bridge",
				Name:      "dhcp1",
			},
			full: &ResourceDHCPServer{
				Interface: "test-bridge",
				Disabled:  false,
				Name:      "dhcp1",
			},
		},
		&testResource{
			min: &ResourceDNSStaticRecord{
				Address: "169.254.169.254",
				Name:    "host.example.tld",
			},
			full: &ResourceDNSStaticRecord{
				Address:  "169.254.169.254",
				Comment:  "test dns record",
				Disabled: false,
				Name:     "host.example.tld",
				TTL:      "1w",
			},
		},
		&testResource{
			min: &ResourceDHCPServerOption{
				Code:  66,
				Name:  "next-server",
				Value: "'169.254.169.1'",
			},
			full: &ResourceDHCPServerOption{
				Code:  67,
				Name:  "bootfile",
				Value: "'pxelinux.0'",
			},
		},
		&testResource{
			env: []Resource{
				&ResourceDHCPServerOption{
					Code:  66,
					Name:  "next-server",
					Value: "'169.254.169.1'",
				},
				&ResourceDHCPServerOption{
					Code:  67,
					Name:  "bootfile",
					Value: "'pxelinux.0'",
				},
			},
			min: &ResourceDHCPServerOptionSet{
				Name:    "PXEClient",
				Options: "next-server",
			},
			full: &ResourceDHCPServerOptionSet{
				Name:    "PXEClient",
				Options: "next-server,bootfile",
			},
		},
		&testResource{
			env: []Resource{
				&ResourceInterfaceBridge{
					Name:     "test-bridge",
					Disabled: false,
					MTU:      1500,
				},
				&ResourceDHCPServer{
					Interface: "test-bridge",
					Disabled:  false,
					Name:      "test-dhcp-server",
				},
				&ResourceDHCPServerOption{
					Code:  66,
					Name:  "next-server",
					Value: "'169.254.169.1'",
				},
			},
			min: &ResourceDHCPServerLease{
				Address:    "169.254.169.254",
				MacAddress: "00:11:22:33:44:55",
				Server:     "test-dhcp-server",
			},
			full: &ResourceDHCPServerLease{
				Address:      "169.254.169.254",
				AddressLists: "none",
				Comment:      "test dhcp server lease",
				ClientID:     "test-machine",
				Disabled:     true,
				DHCPOption:   "next-server",
				MacAddress:   "00:11:22:33:44:55",
				Server:       "test-dhcp-server",
			},
		},
	}
)

func TestCreateResource(t *testing.T) {
	var s *scenario

	c, err := getTestTimeClient()
	if err != nil {
		panic(err)
	}

	defer c.Close()

	conn, connIsStub := c.conn.(*ConnStub)
	if connIsStub {
		s = &scenario{conn: conn}
	}

	for _, r := range resources {
		resource := r.min

		// setup env
		if err, ok := r.setup(c, connIsStub); !ok {
			t.Errorf("unable to setup env before testing: %v", err)
		}

		testName := fmt.Sprintf("when does not exist/%v", reflect.TypeOf(resource))
		t.Run(testName, func(t *testing.T) {
			if connIsStub {
				s.ResourceDoesNotExist()
				s.ResourceCreated()
			}

			if _, err := c.CreateResource(resource); err != nil {
				t.Errorf("expected resource created, got error: %v", err)
			}

		})

		// teardown resource
		if connIsStub {
			s.ResourceExists()
			s.ResourceDeleted()
		}
		if err, ok := c.DeleteResource(resource); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}
		// teardown env
		if err, ok := r.teardown(c, connIsStub); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}
	}

	for _, r := range resources {
		resource := r.min

		// setup env
		if err, ok := r.setup(c, connIsStub); !ok {
			t.Errorf("unable to setup env before testing: %v", err)
		}
		// setup resource: create resource before creating
		if connIsStub {
			s.ResourceDoesNotExist()
			s.ResourceCreated()
		}
		if _, err := c.CreateResource(resource); err != nil {
			t.Errorf("unable to setup env before testing: %v", err)
		}

		testName := fmt.Sprintf("when exists/%v", reflect.TypeOf(resource))
		t.Run(testName, func(t *testing.T) {
			if connIsStub {
				s.ResourceExists()
			}

			if _, err := c.CreateResource(resource); err == nil {
				t.Errorf("expected error, got nil")
			}

		})

		// teardown resource
		if connIsStub {
			s.ResourceExists()
			s.ResourceDeleted()
		}
		if err, ok := c.DeleteResource(resource); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}
		// teardown env
		if err, ok := r.teardown(c, connIsStub); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}
	}

}

func TestReadResource(t *testing.T) {
	var s *scenario

	c, err := getTestTimeClient()
	if err != nil {
		panic(err)
	}

	defer c.Close()

	conn, connIsStub := c.conn.(*ConnStub)
	if connIsStub {
		s = &scenario{conn: conn}
	}

	for _, r := range resources {
		resource := r.min

		testName := fmt.Sprintf("when does not exist/%v", reflect.TypeOf(resource))
		t.Run(testName, func(t *testing.T) {
			if connIsStub {
				s.ResourceDoesNotExist()
			}

			d, err := c.ReadResource(resource)
			if err == nil {
				t.Errorf("expected error, got resource")
			}

			if d != nil {
				t.Errorf("did not expect to receive resource, got resource")
			}
		})
	}

	for _, r := range resources {
		resource := r.min

		// setup env
		if err, ok := r.setup(c, connIsStub); !ok {
			t.Errorf("unable to setup env before testing: %v", err)
		}
		// setup resource: create resource before reading
		if connIsStub {
			s.ResourceDoesNotExist()
			s.ResourceCreated()
		}
		if _, err := c.CreateResource(resource); err != nil {
			t.Errorf("setup error: %v", err)
		}

		testName := fmt.Sprintf("when exists/%v", reflect.TypeOf(resource))
		t.Run(testName, func(t *testing.T) {
			if connIsStub {
				s.ResourceExists()
			}

			d, err := c.ReadResource(resource)
			if err != nil {
				t.Errorf("expected to get resource, got error: %v", err)
			}
			if d == nil {
				t.Errorf("expected to get resource, got nil")
			}

		})

		// teardown resource
		if connIsStub {
			s.ResourceExists()
			s.ResourceDeleted()
		}
		if err, ok := c.DeleteResource(resource); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}
		// teardown env
		if err, ok := r.teardown(c, connIsStub); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}

	}

}

func TestUpdateResource(t *testing.T) {
	var s *scenario

	c, err := getTestTimeClient()
	if err != nil {
		panic(err)
	}

	defer c.Close()

	conn, connIsStub := c.conn.(*ConnStub)
	if connIsStub {
		s = &scenario{conn: conn}
	}

	for _, r := range resources {
		o := r.min
		n := r.full

		testName := fmt.Sprintf("when does not exist/%v", reflect.TypeOf(o))
		t.Run(testName, func(t *testing.T) {
			if connIsStub {
				s.ResourceDoesNotExist()
			}

			if _, ok := c.UpdateResource(o, n); ok {
				t.Errorf("expected to receive err, got ok")
			}
		})
	}

	for _, r := range resources {
		o := r.min
		n := r.full

		// setup env
		if err, ok := r.setup(c, connIsStub); !ok {
			t.Errorf("unable to setup env before testing: %v", err)
		}
		// setup resource: create resource before updating
		if connIsStub {
			s.ResourceDoesNotExist()
			s.ResourceCreated()
		}
		if _, err := c.CreateResource(o); err != nil {
			t.Errorf("unable to setup env before testing: %v", err)
		}

		testName := fmt.Sprintf("when exists/%v", reflect.TypeOf(o))
		t.Run(testName, func(t *testing.T) {

			if connIsStub {
				s.ResourceExists()
				s.ResourceUpdated()
			}

			if err, ok := c.UpdateResource(o, n); !ok {
				t.Errorf("expected to receive ok, got error: %v", err)
			}

		})

		// teardown resource
		if connIsStub {
			s.ResourceExists()
			s.ResourceDeleted()
		}
		if err, ok := c.DeleteResource(n); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}
		// teardown env
		if err, ok := r.teardown(c, connIsStub); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}
	}
}

func TestDeleteResource(t *testing.T) {
	var s *scenario

	c, err := getTestTimeClient()
	if err != nil {
		panic(err)
	}

	defer c.Close()

	conn, connIsStub := c.conn.(*ConnStub)
	if connIsStub {
		s = &scenario{conn: conn}
	}

	for _, r := range resources {
		resource := r.min

		testName := fmt.Sprintf("when does not exist/%v", reflect.TypeOf(resource))
		t.Run(testName, func(t *testing.T) {
			if connIsStub {
				s.ResourceDoesNotExist()
			}

			if _, ok := c.DeleteResource(resource); ok {
				t.Errorf("expected error, got ok")
			}
		})
	}

	for _, r := range resources {
		resource := r.min

		// setup env
		if err, ok := r.setup(c, connIsStub); !ok {
			t.Errorf("unable to setup env before testing: %v", err)
		}
		// setup resource: create resource before deleting
		if connIsStub {
			s.ResourceDoesNotExist()
			s.ResourceCreated()
		}
		if _, err := c.CreateResource(resource); err != nil {
			t.Errorf("unable to setup env before testing: %v", err)
		}

		testName := fmt.Sprintf("when exists/%v", reflect.TypeOf(resource))
		t.Run(testName, func(t *testing.T) {
			if connIsStub {
				s.ResourceExists()
				s.ResourceDeleted()
			}

			if err, ok := c.DeleteResource(resource); !ok {
				t.Errorf("expected resource successfully removed, got error: %v", err)
			}
		})

		// teardown env
		if err, ok := r.teardown(c, connIsStub); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}

	}

}

func TestCheckResourceExists(t *testing.T) {
	var s *scenario

	c, err := getTestTimeClient()
	if err != nil {
		panic(err)
	}

	defer c.Close()

	conn, connIsStub := c.conn.(*ConnStub)
	if connIsStub {
		s = &scenario{conn: conn}
	}

	for _, r := range resources {
		resource := r.min

		testName := fmt.Sprintf("when does not exist/%v", reflect.TypeOf(resource))
		t.Run(testName, func(t *testing.T) {
			if connIsStub {
				s.ResourceDoesNotExist()
			}

			err, ok := c.CheckResourceExists(resource)
			if err != nil {
				t.Errorf("expected !ok, got error: %v", err)
			}
			if ok {
				t.Errorf("expected !ok, got ok")
			}
		})
	}

	for _, r := range resources {
		resource := r.min

		// setup env
		if err, ok := r.setup(c, connIsStub); !ok {
			t.Errorf("unable to setup env before testing: %v", err)
		}
		// setup resource: create resource before checking
		if connIsStub {
			s.ResourceDoesNotExist()
			s.ResourceCreated()
		}
		if _, err := c.CreateResource(resource); err != nil {
			t.Errorf("unable to setup env before testing: %v", err)
		}

		testName := fmt.Sprintf("when exists/%v", reflect.TypeOf(resource))
		t.Run(testName, func(t *testing.T) {
			if connIsStub {
				s.ResourceExists()
			}

			err, ok := c.CheckResourceExists(resource)
			if err != nil {
				t.Errorf("expected ok, got error: %v", err)
			}
			if !ok {
				t.Errorf("expected ok, got !ok")
			}
		})

		// teardown
		if connIsStub {
			s.ResourceExists()
			s.ResourceDeleted()
		}
		if err, ok := c.DeleteResource(resource); !ok {
			t.Errorf("teardown error: %v", err)
		}
		// teardown env
		if err, ok := r.teardown(c, connIsStub); !ok {
			t.Errorf("unable to teardown env after testing: %v", err)
		}
	}
}
