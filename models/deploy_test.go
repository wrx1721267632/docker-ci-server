package models

import (
	"testing"
	"time"
)

func TestDeployAdd(t *testing.T) {
	InitAllInTest()
	deploy := &Deploy{ServiceId: 1, AccountId: 1, DeployStart: time.Now().Unix(), DeployEnd: time.Now().Unix()}
	if _, err := deploy.Add(deploy); err != nil {
		t.Error("Add() failed.Error:", err)
	}
}

func TestDeployUpdate(t *testing.T) {
	InitAllInTest()
	deploy := &Deploy{Id: 1, ServiceId: 1, AccountId: 2}
	if err := deploy.Update(deploy); err != nil {
		t.Error("Update() failed.Error:", err)
	}
}

func TestDeployRemove(t *testing.T) {
	InitAllInTest()
	var deploy Deploy
	if err := deploy.Remove(1); err != nil {
		t.Error("Remove() failed.Error:", err)
	}
}

func TestDeployGetById(t *testing.T) {
	InitAllInTest()
	deploy := &Deploy{ServiceId: 1, AccountId: 1, DeployStart: time.Now().Unix(), DeployEnd: time.Now().Unix()}
	deploy.Add(deploy)

	getDeploy, err := deploy.GetById(deploy.Id)
	if err != nil {
		t.Error("GetById() failed.Error:", err)
	}

	if *getDeploy != *deploy {
		t.Error("GetById() failed", getDeploy, "!=", deploy)
	}
}

func TestDeployGetByServiceId(t *testing.T) {
	InitAllInTest()
	deploy, err := Deploy{}.GetByServiceId(1)
	if err != nil {
		t.Error("GetByDeployId() failed.Error:", err)
	}
	t.Log(deploy)
}
