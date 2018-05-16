package routerosclient

import (
	"testing"
)

var lease = &DHCPLease{
	Address:    "127.0.0.1",
	MacAddress: "11:22:33:44:55:66",
	Server:     "dhcp1",
}

func (s *Scenario) DHCPLeaseCreated() {
	s.conn.buildReply(nil, map[string]string{"ret": "*1"})
}

func (s *Scenario) DHCPLeaseExists() {
	s.conn.buildReply([]map[string]string{{".id": "*1"}}, nil)
}

func (s *Scenario) DHCPLeaseDeleted() {
	s.conn.buildReply(nil, nil)
}

func (s *Scenario) DHCPLeaseDoesNotExist() {
	s.conn.buildReply(nil, nil)
}

func (s *Scenario) DHCPLeaseUpdated() {
	s.conn.buildReply(nil, nil)
}

var (
	client     *Client
	err        error
	conn       *ConnStub
	isConnStub bool
	s          *Scenario
)

func init() {

	if client, err = getTestClient(); err != nil {
		panic(err)
	}

	defer client.Close()

	conn, isConnStub = client.conn.(*ConnStub)
	s = &Scenario{conn: conn}
}

func TestCreateDHCPLease(t *testing.T) {

	teardown := func() {
		if isConnStub {
			s.DHCPLeaseExists()
			s.DHCPLeaseDeleted()
		}
		if err, ok := client.DeleteDHCPLease(lease); !ok {
			t.Errorf("teardown error: %v", err)
		}
	}

	t.Run("when does not exist", func(t *testing.T) {
		if isConnStub {
			s.DHCPLeaseDoesNotExist()
			s.DHCPLeaseCreated()
		}

		if _, err := client.CreateDHCPLease(lease); err != nil {
			t.Errorf("expected lease created, got error: %v", err)
		}
	})

	t.Run("when exists", func(t *testing.T) {
		if isConnStub {
			s.DHCPLeaseExists()
		}

		if _, err := client.CreateDHCPLease(lease); err == nil {
			t.Errorf("expected error, got ok")
		}

	})

	teardown()
}

func TestReadDHCPLease(t *testing.T) {
	var client *Client
	var err error

	if client, err = getTestClient(); err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	conn, isConnStub := client.conn.(*ConnStub)
	s := &Scenario{conn: conn}

	setup := func() {
		if isConnStub {
			s.DHCPLeaseDoesNotExist()
			s.DHCPLeaseCreated()
		}

		if _, err := client.CreateDHCPLease(lease); err != nil {
			t.Errorf("setup error: %v", err)
		}
	}

	teardown := func() {
		if isConnStub {
			s.DHCPLeaseExists()
			s.DHCPLeaseDeleted()
		}
		if err, ok := client.DeleteDHCPLease(lease); !ok {
			t.Errorf("teardown error: %v", err)
		}
	}

	t.Run("when does not exist", func(t *testing.T) {
		if isConnStub {
			s.DHCPLeaseDoesNotExist()
		}

		d, err := client.ReadDHCPLease(lease)
		if err == nil {
			t.Errorf("expected error, got lease")
		}

		if d != nil {
			t.Errorf("did not expect to receive lease, got lease")
		}
	})

	setup()

	t.Run("when exists", func(t *testing.T) {
		if isConnStub {
			s.DHCPLeaseExists()
		}

		d, err := client.ReadDHCPLease(lease)
		if err != nil {
			t.Errorf("expected to get lease, got error: %v", err)
		}
		if d == nil {
			t.Errorf("expected to get lease, got nil")
		}

	})

	teardown()
}

func TestUpdateDHCPLease(t *testing.T) {

	n := &DHCPLease{
		Address:    "169.254.169.254",
		MacAddress: "00:11:22:33:44:55",
		Server:     "dhcp1",
		Disabled:   true,
	}

	setup := func() {
		if isConnStub {
			s.DHCPLeaseDoesNotExist()
			s.DHCPLeaseCreated()
		}
		if _, err := client.CreateDHCPLease(lease); err != nil {
			t.Errorf("setup error: %v", err)
		}
	}

	teardown := func() {
		if isConnStub {
			s.DHCPLeaseExists()
			s.DHCPLeaseDeleted()
		}
		if err, ok := client.DeleteDHCPLease(n); !ok {
			t.Errorf("teardown error: %v", err)
		}
	}

	t.Run("when exists", func(t *testing.T) {
		setup()

		if isConnStub {
			s.DHCPLeaseExists()
			s.DHCPLeaseUpdated()
		}

		if err, ok := client.UpdateDHCPLease(lease, n); !ok {
			t.Errorf("expected to receive ok, got error: %v", err)
		}

		teardown()
	})

	t.Run("when does not exist", func(t *testing.T) {
		if isConnStub {
			s.DHCPLeaseDoesNotExist()
		}

		if _, ok := client.UpdateDHCPLease(lease, n); ok {
			t.Errorf("expected to receive err, got ok")
		}
	})
}

func TestDeleteDHCPLease(t *testing.T) {

	setup := func() {
		if isConnStub {
			s.DHCPLeaseDoesNotExist()
			s.DHCPLeaseCreated()
		}

		if _, err := client.CreateDHCPLease(lease); err != nil {
			t.Errorf("setup error: %v", err)
		}
	}

	t.Run("when does not exist", func(t *testing.T) {
		if isConnStub {
			s.DHCPLeaseDoesNotExist()
		}

		if _, ok := client.DeleteDHCPLease(lease); ok {
			t.Errorf("expected error, got ok")
		}
	})

	setup()

	t.Run("when exists", func(t *testing.T) {
		if isConnStub {
			s.DHCPLeaseExists()
			s.DHCPLeaseDeleted()
		}

		if err, ok := client.DeleteDHCPLease(lease); !ok {
			t.Errorf("expected lease successfully removed, got error: %v", err)
		}
	})
}

func TestCheckDHCPLeaseExists(t *testing.T) {

	setup := func() {
		if isConnStub {
			s.DHCPLeaseDoesNotExist()
			s.DHCPLeaseCreated()
		}

		if _, err := client.CreateDHCPLease(lease); err != nil {
			t.Errorf("setup error: %v", err)
		}
	}

	teardown := func() {
		if isConnStub {
			s.DHCPLeaseExists()
			s.DHCPLeaseDeleted()
		}
		if err, ok := client.DeleteDHCPLease(lease); !ok {
			t.Errorf("teardown error: %v", err)
		}
	}

	t.Run("when does not exist", func(t *testing.T) {
		if isConnStub {
			s.DHCPLeaseDoesNotExist()
		}

		err, ok := client.CheckDHCPLeaseExists(lease)
		if err != nil {
			t.Errorf("expected !ok, got error: %v", err)
		}
		if ok {
			t.Errorf("expected !ok, got ok")
		}
	})

	setup()

	t.Run("when exists", func(t *testing.T) {
		if isConnStub {
			s.DHCPLeaseExists()
		}

		err, ok := client.CheckDHCPLeaseExists(lease)
		if err != nil {
			t.Errorf("expected ok, got error: %v", err)
		}
		if !ok {
			t.Errorf("expected ok, got !ok")
		}
	})

	teardown()
}
