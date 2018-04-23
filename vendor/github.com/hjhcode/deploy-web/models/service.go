package models

import (
	"time"

	. "github.com/hjhcode/deploy-web/common/store"
)

type Service struct {
	Id              int64
	AccountId       int64
	ServiceName     string
	ServiceDescribe string
	HostList        string
	MirrorList      string
	DockerConfig    string
	CreateDate      time.Time
	UpdateDate      time.Time
	ServiceMember   string
}

//增加
func (this Service) Add(service *Service) (int64, error) {
	_, err := OrmWeb.Insert(service)
	if err != nil {
		return 0, err
	}
	return service.Id, nil
}

//删除
func (this Service) Remove(id int64) error {
	service := new(Service)
	_, err := OrmWeb.Id(id).Delete(service)
	return err
}

//修改
func (this Service) Update(service *Service) error {
	_, err := OrmWeb.Id(service.Id).Update(service)
	return err
}

//查询(根据服务id查询)
func (this Service) GetById(id int64) (*Service, error) {
	service := new(Service)
	has, err := OrmWeb.Id(id).Get(service)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return service, nil
}

//查询(根据服务名模糊查询）
func (this Service) QueryByServiceName(serviceName string) ([]*Service, error) {
	serviceList := make([]*Service, 0)
	err := OrmWeb.Where("service_name like ?", "%"+serviceName+"%").Find(&serviceList)
	if err != nil {
		return nil, err
	}
	return serviceList, nil
}

//查询(根据服务创建者查询）
func (this Service) QueryByAccountId(accountId int64) ([]*Service, error) {
	serviceList := make([]*Service, 0)
	err := OrmWeb.Where("account_id = ?", accountId).Find(&serviceList)
	if err != nil {
		return nil, err
	}
	return serviceList, nil
}

//查询(查询所有服务)
func (this Service) QueryAllService() ([]*Service, error) {
	serviceList := make([]*Service, 0)
	err := OrmWeb.Find(&serviceList)
	if err != nil {
		return nil, err
	}
	return serviceList, nil
}
