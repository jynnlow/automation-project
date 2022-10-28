package model

import (
	"automation-project/dto"
	"time"

	"gorm.io/gorm"
)

type MapCRUDInterface interface {
	Create(*Map) (*Map, error)
	List() (*[]Map, error)
	GetByKey(string) (*Map, error)
	Update(*Map) (*Map, error)
	Delete(string) (*Map, error)
}

type Map struct {
	ID         uint   `gorm:"primaryKey"`
	Key        string `gorm:"unique"`
	Value      int
	Created_at int
	Updated_at int
}

type MapCRUD struct {
	DB *gorm.DB
}

func (m *MapCRUD) Create(mapReq *dto.MapReq) (*Map, error) {
	newMap := &Map{
		Key:        mapReq.Key,
		Value:      mapReq.Value,
		Created_at: int(time.Now().Unix()),
		Updated_at: 0,
	}
	err := m.DB.Create(newMap).Error
	if err != nil {
		return nil, err
	}
	return newMap, nil
}

func (m *MapCRUD) List() ([]*Map, error) {
	var mapList []*Map
	err := m.DB.Find(&mapList).Error
	if err != nil {
		return nil, err
	}
	return mapList, nil
}

func (m *MapCRUD) GetByKey(key string) (*Map, error) {
	foundMap := &Map{}
	err := m.DB.Where("`key` = ?", key).First(foundMap).Error
	if err != nil {
		return nil, err
	}
	return foundMap, nil
}

func (m *MapCRUD) Update(mapReq *dto.MapReq) (*Map, error) {
	foundMap, err := m.GetByKey(mapReq.Key)
	if err != nil {
		return nil, err
	}

	foundMap.Value = mapReq.Value
	foundMap.Updated_at = int(time.Now().Unix())
	err = m.DB.Save(foundMap).Error
	if err != nil {
		return nil, err
	}
	return foundMap, nil
}

func (m *MapCRUD) Delete(key string) (*Map, error) {
	foundMap, err := m.GetByKey(key)
	if err != nil {
		return nil, err
	}
	//delete only if the user exists
	//permanently deleted with Unscoped().Delete()
	err = m.DB.Where("`key` = ?", key).Unscoped().Delete(foundMap).Error
	if err != nil {
		return nil, err
	}
	return foundMap, nil
}
