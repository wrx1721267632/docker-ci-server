package models

import (
	"testing"
	"time"
)

func TestConstructRecordAdd(t *testing.T) {
	InitAllInTest()
	constructRecord := &ConstructRecord{AccountId: 1, ProjectId: 1, ConstructStart: time.Now().Unix(),
		ConstructEnd: time.Now().Unix()}
	if _, err := constructRecord.Add(constructRecord); err != nil {
		t.Error("Add() failed.Error:", err)
	}
}

func TestConstructRecordUpdate(t *testing.T) {
	InitAllInTest()
	constructRecord := &ConstructRecord{Id: 2, AccountId: 1, ProjectId: 3}
	if err := constructRecord.Update(constructRecord); err != nil {
		t.Error("Update() failed.Error:", err)
	}
}

func TestConstructRecordRemove(t *testing.T) {
	InitAllInTest()
	var constructRecord ConstructRecord
	if err := constructRecord.Remove(1); err != nil {
		t.Error("Remove() failed.Error:", err)
	}
}

func TestConstructRecordGetById(t *testing.T) {
	InitAllInTest()
	constructRecord := &ConstructRecord{AccountId: 1, ProjectId: 1, ConstructStart: time.Now().Unix(),
		ConstructEnd: time.Now().Unix()}
	constructRecord.Add(constructRecord)

	getConstructRecord, err := constructRecord.GetById(constructRecord.Id)
	if err != nil {
		t.Error("GetById() failed.Error:", err)
	}

	if *getConstructRecord != *constructRecord {
		t.Error("GetById() failed", getConstructRecord, "!=", constructRecord)
	}
}

func TestConstructRecordGetByProjectId(t *testing.T) {
	InitAllInTest()
	record, err := ConstructRecord{}.GetByProjectId(1)
	if err != nil {
		t.Error("GetByProjectId() failed.Error:", err)
	}
	t.Log(record)
}
