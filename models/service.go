package models

import (
	. "github.com/wrxcode/deploy-server/common/store"
)

type Service struct {
	Id              int64  `json:"id" form:"id"`
	AccountId       int64  `json:"account_id" form:"account_id"`
	ServiceName     string `json:"service_name" form:"service_name"`
	ServiceDescribe string `json:"service_describe" form:"service_describe"`
	HostList        string `json:"host_list" form:"host_list"`
	MirrorList      int64  `json:"mirror_list" form:"mirror_list"`
	DockerConfig    string `json:"docker_config" form:"docker_config"`
	CreateDate      int64  `json:"create_date" form:"create_date"`
	UpdateDate      int64  `json:"update_date" form:"update_date"`
	ServiceMember   string `json:"service_member" form:"service_member"`
	IsDel           int    `json:"is_del" form:"is_del"`
	ServiceStatu    int    `json:"service_statu" form:"service_statu"`
}

type ServiceData struct {
	Stage []ServiceStage `json:"stage"`
}

type ServiceStage struct {
	Machine []ServiceMachine `json:"machine"`
}

type ServiceMachine struct {
	Id               int64  `json:"id"`
	Container_status string `json:"container_status"`
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

//查询(根据服务名精确查询)
func (this Service) GetByServiceName(serviceName string) (*Service, error) {
	service := new(Service)
	has, err := OrmWeb.Where("service_name=?", serviceName).Get(service)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return service, nil
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
	err := OrmWeb.Where("is_del != ?", 1).Find(&serviceList)
	if err != nil {
		return nil, err
	}
	return serviceList, nil
}

//查询(分页查询所有服务）
func (this Service) QueryAllServiceByPage(size int, start int) ([]*Service, error) {
	serviceList := make([]*Service, 0)
	err := OrmWeb.Where("is_del != ?", 1).Limit(size, start).Find(&serviceList)
	if err != nil {
		return nil, err
	}

	return serviceList, nil
}

func (this Service) CountAllService() (int64, error) {
	sum, err := OrmWeb.Count(&Service{})
	if err != nil {
		return 0, nil
	}

	return sum, nil
}

//查询(根据服务名查询）
func (this Service) QueryServiceBySearch(serviceName string, service *Service) ([]*Service, error) {
	serviceList := make([]*Service, 0)
	err := OrmWeb.Where("service_name like ?", "%"+serviceName+"%").Find(&serviceList, service)
	if err != nil {
		return nil, err
	}

	return serviceList, nil
}

func (this Service) CountServiceBySearch(serviceName string, service *Service) (int64, error) {
	sum, err := OrmWeb.Where("service_name like ?", "%"+serviceName+"%").Count(service)
	if err != nil {
		return 0, err
	}

	return sum, nil
}
