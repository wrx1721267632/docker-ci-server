package models

import (
	. "github.com/hjhcode/deploy-web/common/store"
)

type Deploy struct {
	Id           int64
	ServiceId    int64
	AccountId    int64
	DeployStart  int64
	DeployEnd    int64
	HostList     string
	MirrorList   string
	DockerConfig string
	DeployStatu  int
	DeployLog    string
}

//增加
func (this Deploy) Add(deploy *Deploy) (int64, error) {
	_, err := OrmWeb.Insert(deploy)
	if err != nil {
		return 0, err
	}
	return deploy.Id, nil
}

//删除
func (this Deploy) Remove(id int64) error {
	deploy := new(Deploy)
	_, err := OrmWeb.Id(id).Delete(deploy)
	return err
}

//修改
func (this Deploy) Update(deploy *Deploy) error {
	_, err := OrmWeb.Id(deploy.Id).Update(deploy)
	return err
}

//查询(根据部署id查询)
func (this Deploy) GetById(id int64) (*Deploy, error) {
	deploy := new(Deploy)
	has, err := OrmWeb.Id(id).Get(deploy)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return deploy, nil
}

//查询(根据服务id查询）
func (this Deploy) QueryByServiceId(serviceId int64) ([]*Deploy, error) {
	deployList := make([]*Deploy, 0)
	err := OrmWeb.Where("service_id = ?", serviceId).Find(&deployList)
	if err != nil {
		return nil, err
	}
	return deployList, nil
}

//查询(根据用户id查询）
func (this Deploy) QueryByAccountId(accountId int64) ([]*Deploy, error) {
	deployList := make([]*Deploy, 0)
	err := OrmWeb.Where("account_id=?", accountId).Find(&deployList)
	if err != nil {
		return nil, err
	}
	return deployList, err
}

//查询(所有部署）
func (this Deploy) QueryAllDeploy() ([]*Deploy, error) {
	deployList := make([]*Deploy, 0)
	err := OrmWeb.Find(&deployList)
	if err != nil {
		return nil, err
	}
	return deployList, err
}
