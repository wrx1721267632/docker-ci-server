package models

import (
	. "github.com/wrxcode/deploy-server/common/store"
)

type Mirror struct {
	Id             int64  `json:"id" form:"id"`
	MirrorName     string `json:"mirror_name" form:"mirror_name"`
	MirrorVersion  string `json:"mirror_version" form:"mirror_version"`
	MirrorDescribe string `json:"mirror_describe" form:"mirror_describe"`
}

//增加
func (this Mirror) Add(mirror *Mirror) (int64, error) {
	_, err := OrmWeb.Insert(mirror)
	if err != nil {
		return 0, err
	}
	return mirror.Id, nil
}

//删除
func (this Mirror) Remove(id int64) error {
	mirror := new(Mirror)
	_, err := OrmWeb.Id(id).Delete(mirror)
	return err
}

//修改
func (this Mirror) Update(mirror *Mirror) error {
	_, err := OrmWeb.Id(mirror.Id).Update(mirror)
	return err
}

//查询(根据镜像id查询）
func (this Mirror) GetById(id int64) (*Mirror, error) {
	mirror := new(Mirror)
	has, err := OrmWeb.Id(id).Get(mirror)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return mirror, nil
}

//查询(根据镜像名字查询）
func (this Mirror) QueryByMirrorName(mirrorName string) ([]*Mirror, error) {
	mirrorList := make([]*Mirror, 0)
	err := OrmWeb.Find(&mirrorList)
	if err != nil {
		return nil, err
	}
	return mirrorList, nil
}
