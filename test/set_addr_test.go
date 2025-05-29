package test

import (
	"raspberrypi/utils/network"
	"testing"
)

func TestSetAddr(t *testing.T) {
	err := network.SetIPv4Addr("eth0", "192.168.0.1/30")
	if err != nil {
		t.Errorf("failed to set address: %v", err)
	} else {
		t.Log("success")
	}
}
