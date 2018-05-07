package models

import (
	. "github.com/wrxcode/deploy-server/common/store"
)

type Deploy struct {
	Id           int64  `json:"id" form:"id"`
	ServiceId    int64  `json:"service_id" form:"service_id"`
	AccountId    int64  `json:"account_id" form:"account_id"`
	DeployStart  int64  `json:"deploy_start" form:"deploy_start"`
	DeployEnd    int64  `json:"deploy_end" form:"deploy_end"`
	HostList     string `json:"host_list" form:"host_list"`
	MirrorList   int64  `json:"mirror_list" form:"mirror_list"`
	DockerConfig string `json:"docker_config" form:"docker_config"`
	DeployStatu  int    `json:"deploy_statu" form:"deploy_statu"`
	DeployLog    string `json:"deploy_log" form:"deploy_log"`
}

type DeployData struct {
	Stage           []DeployStage `json:"stage"`
	Stage_num       int           `json:"stage_num"`
	Progress_status int           `json:"progress_status"`
}

type DeployStage struct {
	Stage_status int             `json:"stage_status"`
	Machine      []DeployMachine `json:"machine"`
}

type DeployMachine struct {
	Id             int64  `json:"id"`
	Machine_status int    `json:"machine_status"`
	Step           string `json:"step"`
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
	err := OrmWeb.Desc("id").Find(&deployList)
	if err != nil {
		return nil, err
	}
	return deployList, err
}

//查询(分页查询所有记录）
func (this Deploy) QueryAllDeployByPage(size int, start int) ([]*Deploy, error) {
	deployList := make([]*Deploy, 0)
	err := OrmWeb.Limit(size, start).Find(&deployList)
	if err != nil {
		return nil, err
	}

	return deployList, nil
}

func (this Deploy) CountAllDeployByPage() (int64, error) {
	sum, err := OrmWeb.Count(&Deploy{})
	if err != nil {
		return 0, nil
	}

	return sum, nil
}

func (this Deploy) GetByServiceId(serviceId int64) (*Deploy, error) {
	deploy := &Deploy{ServiceId: serviceId}
	has, err := OrmWeb.Desc("id").Get(deploy)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return deploy, nil
}

func (this Deploy) GetDeployBackData(serviceId int64) (*Deploy, error) {
	deploy := &Deploy{ServiceId: serviceId, DeployStatu: 3}
	has, err := OrmWeb.Desc("id").Get(deploy)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return deploy, nil
}
