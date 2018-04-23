package models

import (
	"time"

	. "github.com/hjhcode/deploy-web/common/store"
)

type ConstructRecord struct {
	Id             int64
	AccountId      int64
	ProjectId      int64
	MirrorId       int64
	ConstructStart time.Time
	ConstructEnd   time.Time
	ConstructStatu int
	ConstructLog   string
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
	err := OrmWeb.Find(&constructRecordList)
	if err != nil {
		return nil, err
	}
	return constructRecordList, nil
}
