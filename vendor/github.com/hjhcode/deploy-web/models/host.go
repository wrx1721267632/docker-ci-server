package models

import (
	. "github.com/hjhcode/deploy-web/common/store"
)

type Host struct {
	Id       int64
	HostName string
	Ip       string
}

//增加
func (this Host) Add(host *Host) (int64, error) {
	_, err := OrmWeb.Insert(host)
	if err != nil {
		return 0, err
	}
	return host.Id, nil
}

//删除
func (this Host) Remove(id int64) error {
	host := new(Host)
	_, err := OrmWeb.Id(id).Delete(host)
	return err
}

//修改
func (this Host) Update(host *Host) error {
	_, err := OrmWeb.Id(host.Id).Update(host)
	return err
}

//查询（根据主机id查询）
func (this Host) GetById(id int64) (*Host, error) {
	host := new(Host)
	has, err := OrmWeb.Id(id).Get(host)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return host, nil
}

//查询（根据主机名字查询）
func (this Host) GetByHostName(hostName string) (*Host, error) {
	host := new(Host)
	has, err := OrmWeb.Where("host_name=?", hostName).Get(host)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return host, nil
}

//查询(根据ip查询）
func (this Host) GetByIp(ip string) (*Host, error) {
	host := new(Host)
	has, err := OrmWeb.Where("ip=?", ip).Get(host)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return host, nil
}

//查询(所有数据）
func (this Host) QueryAllHost() ([]*Host, error) {
	hostList := make([]*Host, 0)
	err := OrmWeb.Find(&hostList)
	if err != nil {
		return nil, err
	}
	return hostList, nil
}
