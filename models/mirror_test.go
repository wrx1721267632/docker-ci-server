package models

import (
	"testing"
)

func TestMirrorAdd(t *testing.T) {
	InitAllInTest()
	mirror := &Mirror{MirrorName: "aaaa"}
	if _, err := mirror.Add(mirror); err != nil {
		t.Error("Add() failed.Error:", err)
	}
}

func TestMirrorUpdate(t *testing.T) {
	InitAllInTest()
	mirror := &Mirror{Id: 1, MirrorName: "jjjj"}
	if err := mirror.Update(mirror); err != nil {
		t.Error("Update() failed.Error:", err)
	}
}

func TestMirrorRemove(t *testing.T) {
	InitAllInTest()
	var mirror Mirror
	if err := mirror.Remove(1); err != nil {
		t.Error("Remove() failed.Error:", err)
	}
}

func TestMirrorGetById(t *testing.T) {
	InitAllInTest()
	mirror := &Mirror{MirrorName: "aaaa"}
	mirror.Add(mirror)

	getMirror, err := mirror.GetById(mirror.Id)
	if err != nil {
		t.Error("GetById() failed.Error:", err)
	}

	if *getMirror != *mirror {
		t.Error("GetById() failed", getMirror, "!=", mirror)
	}
}
