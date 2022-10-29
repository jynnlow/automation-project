package model

import (
	"automation-project/dto"

	"gorm.io/gorm"
)

type AuditLogCRUDInterface interface {
	Create(*AuditLog) (*AuditLog, error)
	List() ([]*AuditLog, error)
	GetListByKey(string) ([]*AuditLog, error)
	GetLatestByKey(string) (*AuditLog, error)
	GetByTimstamp(int) (*AuditLog, error)
	UpdateIsLatest(*AuditLog) (*AuditLog, error)
}

type AuditLog struct {
	ID             uint `gorm:"primaryKey"`
	Timestamp      int  `gorm:"unique"`
	Map_id         uint
	Key            string
	Original_value int
	New_value      int
	Action         string
	Is_latest      bool
	User_id        int
}

type AuditLogCRUD struct {
	DB *gorm.DB
}

func (a *AuditLogCRUD) Create(auditReq *dto.AuditLogReq) (*AuditLog, error) {
	newAudit := &AuditLog{
		Timestamp:      auditReq.Timestamp,
		Map_id:         auditReq.Map_id,
		Key:            auditReq.Key,
		Original_value: auditReq.Original_value,
		New_value:      auditReq.New_value,
		Action:         auditReq.Action,
		Is_latest:      auditReq.Is_latest,
		User_id:        auditReq.User_id,
	}
	err := a.DB.Create(newAudit).Error
	if err != nil {
		return nil, err
	}
	return newAudit, nil
}

func (a *AuditLogCRUD) List() ([]*AuditLog, error) {
	var auditList []*AuditLog
	err := a.DB.Find(&auditList).Error
	if err != nil {
		return nil, err
	}
	return auditList, nil
}

func (a *AuditLogCRUD) GetListByKey(key string) ([]*AuditLog, error) {
	var auditList []*AuditLog
	err := a.DB.Where("`key` = ?", key).Find(&auditList).Error
	if err != nil {
		return nil, err
	}
	return auditList, nil
}

func (a *AuditLogCRUD) GetLatestByKey(key string) (*AuditLog, error) {
	foundLog := &AuditLog{}
	err := a.DB.Where("`key` = ?", key).Last(foundLog).Error
	if err != nil {
		return nil, err
	}
	return foundLog, nil
}

func (a *AuditLogCRUD) GetByTimstamp(timestamp int) (*AuditLog, error) {
	foundLog := &AuditLog{}
	err := a.DB.Where("`timestamp` = ?", timestamp).First(foundLog).Error
	if err != nil {
		return nil, err
	}
	return foundLog, nil
}

func (a *AuditLogCRUD) UpdateIsLatest(key string) (*AuditLog, error) {
	foundLog, err := a.GetLatestByKey(key)
	if err != nil {
		return nil, err
	}

	err = a.DB.Model(foundLog).Update("is_latest", false).Error
	if err != nil {
		return nil, err
	}
	return foundLog, nil
}
