package models

import (
	. "github.com/wrxcode/deploy-server/common/store"
)

type ConstructRecord struct {
	Id             int64  `json:"id" form:"id"`
	AccountId      int64  `json:"account_id" form:"account_id"`
	ProjectId      int64  `json:"project_id" form:"project_id"`
	MirrorId       int64  `json:"mirror_id" form:"mirror_id"`
	ConstructStart int64  `json:"construct_start" form:"construct_start"`
	ConstructEnd   int64  `json:"construct_end" form:"construct_end"`
	ConstructStatu int    `json:"construct_statu" form:"construct_statu"`
	ConstructLog   string `json:"construct_log" form:"construct_log"`
}

//增加
func (this ConstructRecord) Add(constructRecord *ConstructRecord) (int64, error) {
	_, err := OrmWeb.Insert(constructRecord)
	if err != nil {
		return 0, err
	}
	return constructRecord.Id, nil
}

//删除
func (this ConstructRecord) Remove(id int64) error {
	constructRecord := new(ConstructRecord)
	_, err := OrmWeb.Id(id).Delete(constructRecord)
	return err
}

//修改
func (this ConstructRecord) Update(constructRecord *ConstructRecord) error {
	_, err := OrmWeb.Id(constructRecord.Id).Update(constructRecord)
	return err
}

//查询（根据构建记录id查询）
func (this ConstructRecord) GetById(id int64) (*ConstructRecord, error) {
	constructRecord := new(ConstructRecord)
	has, err := OrmWeb.Id(id).Get(constructRecord)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return constructRecord, nil

}

//查询（根据工程id查询）
func (this ConstructRecord) QueryByProjectId(projectId int64) ([]*ConstructRecord, error) {
	constructRecordList := make([]*ConstructRecord, 0)
	err := OrmWeb.Where("project_id = ?", projectId).Find(&constructRecordList)
	if err != nil {
		return nil, err
	}
	return constructRecordList, nil
}

//查询（根据用户id查询）
func (this ConstructRecord) QueryByAccountId(accountId int64) ([]*ConstructRecord, error) {
	constructRecordList := make([]*ConstructRecord, 0)
	err := OrmWeb.Where("account_id = ?", accountId).Find(&constructRecordList)
	if err != nil {
		return nil, err
	}
	return constructRecordList, nil
}

//查询(所有数据）
func (this ConstructRecord) QueryAllConstructRecord() ([]*ConstructRecord, error) {
	constructRecordList := make([]*ConstructRecord, 0)
	err := OrmWeb.Desc("id").Find(&constructRecordList)
	if err != nil {
		return nil, err
	}
	return constructRecordList, nil
}

//查询(分页查询所有记录）
func (this ConstructRecord) QueryAllConstructByPage(size int, start int) ([]*ConstructRecord, error) {
	constructList := make([]*ConstructRecord, 0)
	err := OrmWeb.Limit(size, start).Find(&constructList)
	if err != nil {
		return nil, err
	}

	return constructList, nil
}

func (this ConstructRecord) CountAllConstructByPage() (int64, error) {
	sum, err := OrmWeb.Count(&ConstructRecord{})
	if err != nil {
		return 0, nil
	}

	return sum, nil
}

func (this ConstructRecord) GetByProjectId(projectId int64) (*ConstructRecord, error) {
	constructRecord := &ConstructRecord{ProjectId: projectId}
	has, err := OrmWeb.Desc("id").Get(constructRecord)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return constructRecord, nil
}

func (this ConstructRecord) GetByProjectIdAndStatu(projectId int64) (*ConstructRecord, error) {
	constructRecord := &ConstructRecord{ProjectId: projectId, ConstructStatu: 2}
	has, err := OrmWeb.Desc("id").Get(constructRecord)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return constructRecord, nil
}
