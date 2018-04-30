package models

import (
	"testing"
)

func TestHostAdd(t *testing.T) {
	InitAllInTest()
	host := &Host{HostName: "aaaaaa"}
	if _, err := host.Add(host); err != nil {
		t.Error("Add() failed.Error:", err)
	}
}

func TestHostUpdate(t *testing.T) {
	InitAllInTest()
	host := &Host{Id: 1, HostName: "aaxxxx"}
	if err := host.Update(host); err != nil {
		t.Error("Update() failed.Error:", err)
	}
}

func TestHostRemove(t *testing.T) {
	InitAllInTest()
	var host Host
	if err := host.Remove(1); err != nil {
		t.Error("Remove() failed.Error:", err)
	}
}

func TestHostGetById(t *testing.T) {
	InitAllInTest()
	host := &Host{HostName: "aaaaaa"}
	host.Add(host)

	getHost, err := host.GetById(host.Id)
	if err != nil {
		t.Error("GetById() failed.Error:", err)
	}

	if *getHost != *host {
		t.Error("GetById() failed", getHost, "!=", host)
	}
}
