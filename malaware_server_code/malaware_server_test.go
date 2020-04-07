package main

import (
	"../Execute_systems_commands"
	"testing"
)

func TestLocalIP(t *testing.T) {
	localip, err := Execute_systems_commands.RetrieveCurrentIP()
	if err != nil {
		t.Fail()
	}
	if localip != "10.0.2.15" {
		t.Fail()
	}
}
