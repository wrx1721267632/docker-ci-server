package models

import (
	"testing"
)

func TestAccountAdd(t *testing.T) {
	InitAllInTest()

	account := &Account{Name: "fffff@qq.com", Password: "123"}
	if _, err := account.Add(account); err != nil {
		t.Error("Add() failed.Error:", err)
	}
}

func TestAccountUpdate(t *testing.T) {
	InitAllInTest()

	account := &Account{1, "qaqq@qq.com", "88"}
	if err := account.Update(account); err != nil {
		t.Error("Update() failed.Error:", err)
	}
}

func TestAccountRemove(t *testing.T) {
	InitAllInTest()

	var account Account
	if err := account.Remove(4); err != nil {
		t.Error("Remove() failed.Error:", err)
	}
}

func TestAccountGetById(t *testing.T) {
	InitAllInTest()

	account := &Account{Name: "bbb@qq.com", Password: "123"}
	account.Add(account)

	getAccount, err := account.GetById(account.Id)
	if err != nil {
		t.Error("GetById() failed:", err.Error())
	}

	if *getAccount != *account {
		t.Error("GetById() failed:", "%v != %v", account, getAccount)
	}
}

func TestAccountGetByName(t *testing.T) {
	InitAllInTest()

	account := &Account{Name: "bbbs@qq.com", Password: "123"}
	account.Add(account)

	getAccount, err := account.GetByName(account.Name)
	if err != nil {
		t.Error("GetById() failed:", err.Error())
	}

	if *getAccount != *account {
		t.Error("GetById() failed:", account, "!=", getAccount)
	}
}
