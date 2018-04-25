package models

import (
	. "github.com/hjhcode/deploy-web/common/store"
)

type Project struct {
	Id              int64
	AccountId       int64
	ProjectName     string
	ProjectDescribe string
	GitDockerPath   string
	CreateDate      int64
	UpdateDate      int64
	IsDel           int
	ProjectMember   string
}

//增加
func (this Project) Add(project *Project) (int64, error) {
	_, err := OrmWeb.Insert(project)
	if err != nil {
		return 0, err
	}
	return project.Id, nil
}

//删除
func (this Project) Remove(id int64) error {
	project := new(Project)
	_, err := OrmWeb.Id(id).Delete(project)
	return err
}

//修改
func (this Project) Update(project *Project) error {
	_, err := OrmWeb.Id(project.Id).Update(project)
	return err
}

//查询（根据工程id查询）
func (this Project) GetById(id int64) (*Project, error) {
	project := new(Project)
	has, err := OrmWeb.Id(id).Get(project)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return project, nil
}

//查询（根据工程名模糊查询）
func (this Project) QueryByProjectName(projectName string) ([]*Project, error) {
	projectList := make([]*Project, 0)
	err := OrmWeb.Where("project_name like ?", "%"+projectName+"%").Find(&projectList)
	if err != nil {
		return nil, err
	}
	return projectList, nil
}

//查询(根据工程名精确查询）
func (this Project) GetByProjectName(projectName string) (*Project, error) {
	project := new(Project)
	has, err := OrmWeb.Where("project_name=?", projectName).Get(project)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return project, nil
}

//查询(根据创建者查询）
func (this Project) QueryByAccountId(accountId int64) ([]*Project, error) {
	projectList := make([]*Project, 0)
	err := OrmWeb.Where("account_id = ?", accountId).Find(&projectList)
	if err != nil {
		return nil, err
	}
	return projectList, nil
}

//查询(查询所有工程）
func (this Project) QueryAllProject() ([]*Project, error) {
	projectList := make([]*Project, 0)
	err := OrmWeb.Find(&projectList)
	if err != nil {
		return nil, err
	}
	return projectList, nil
}

//查询(分页查询所有工程）
func (this Project) QueryAllProjectByPage(size int, start int) ([]*Project, error) {
	projectList := make([]*Project, 0)
	err := OrmWeb.Where("is_del != ?", 1).Limit(size, start).Find(&projectList)
	if err != nil {
		return nil, err
	}
	return projectList, nil
}

func (this Project) CountAllProjectByPage() (int64, error) {
	sum, err := OrmWeb.Where("is_del != ?", 1).Count(&Project{})
	if err != nil {
		return 0, err
	}
	return sum, nil
}

//查询(根据工程名查询）
func (this Project) QueryProjectBySearch(projectName string, project *Project, size int, start int) ([]*Project, error) {
	projectList := make([]*Project, 0)
	err := OrmWeb.Where("project_name like ?", "%"+projectName+"%").
		Limit(size, start).Find(&projectList, project)
	if err != nil {
		return nil, err
	}

	return projectList, nil
}

func (this Project) CountBySearch(projectName string, project *Project) (int64, error) {
	sum, err := OrmWeb.Where("project_name like ?", "%"+projectName+"%").Count(project)
	if err != nil {
		return 0, err
	}

	return sum, nil
}
