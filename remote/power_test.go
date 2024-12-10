package remote

import "testing"

func TestSetPower(t *testing.T) {
	control, adapter := newControlTester()

	control.SetPower(true)
	if adapter.String() != "-p.1\r" {
		t.Fail()
	}
	adapter.Reset()

	control.SetPower(false)
	if adapter.String() != "-p.0\r" {
		t.Fail()
	}
	adapter.Reset()
}

func TestTogglePower(t *testing.T) {
	control, adapter := newControlTester()

	control.TogglePower()
	if adapter.String() != "-p.t\r" {
		t.Fail()
	}
	adapter.Reset()
}
