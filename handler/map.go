package handler

import (
	"automation-project/dto"
	"automation-project/helper"
	"automation-project/model"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"
)

type MapHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	GetByKey(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Revert(w http.ResponseWriter, r *http.Request)
}

type MapHandler struct {
	MapModel   model.MapCRUDInterface
	AuditModel model.AuditLogCRUDInterface
}

func (m *MapHandler) Create(w http.ResponseWriter, r *http.Request) {
	mapReq := &dto.MapReq{}
	err := json.NewDecoder(r.Body).Decode(mapReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusBadRequest, err)
		return
	}

	//create in main map table
	dbMap, err := m.MapModel.Create(mapReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	auditReq := &dto.AuditLogReq{
		Timestamp:      dbMap.Created_at,
		Map_id:         dbMap.ID,
		Key:            dbMap.Key,
		Original_value: 0,
		New_value:      dbMap.Value,
		Action:         "created",
		Is_latest:      true,
		User_id:        rand.Intn(100-1) + 1,
	}

	//insert into audit log table
	_, err = m.AuditModel.Create(auditReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessRes(w, "Created Successfully", dbMap)
}

func (m *MapHandler) List(w http.ResponseWriter, r *http.Request) {
	dbMapList, err := m.MapModel.List()
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessRes(w, "Retrieved successfully", dbMapList)
}

func (m *MapHandler) GetByKey(w http.ResponseWriter, r *http.Request) {
	param, exist := r.URL.Query()["key"]
	if !exist || len(param[0]) < 1 {
		helper.ErrorRes(w, http.StatusBadRequest, errors.New("param not exist"))
		return
	}

	key := param[0]

	dbMap, err := m.MapModel.GetByKey(key)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessRes(w, "Key found", dbMap)
}

func (m *MapHandler) Update(w http.ResponseWriter, r *http.Request) {
	mapReq := &dto.MapReq{}
	err := json.NewDecoder(r.Body).Decode(mapReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusBadRequest, err)
		return
	}

	//get ori data from map table
	dbOriMap, err := m.MapModel.GetByKey(mapReq.Key)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	//update in the main map table
	dbUpdatedMap, err := m.MapModel.Update(mapReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	auditReq := &dto.AuditLogReq{
		Timestamp:      dbUpdatedMap.Updated_at,
		Map_id:         dbUpdatedMap.ID,
		Key:            dbUpdatedMap.Key,
		Original_value: dbOriMap.Value,
		New_value:      dbUpdatedMap.Value,
		Action:         "updated",
		Is_latest:      true,
		User_id:        generateRandom(),
	}

	//update audit log record is latest = false
	_, err = m.AuditModel.UpdateIsLatest(mapReq.Key)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	//insert update operation into audit log table
	_, err = m.AuditModel.Create(auditReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessRes(w, "Updated successfully", dbUpdatedMap)
}

func (m *MapHandler) Delete(w http.ResponseWriter, r *http.Request) {
	param, exist := r.URL.Query()["key"]
	if !exist || len(param[0]) < 1 {
		helper.ErrorRes(w, http.StatusBadRequest, errors.New("param not exist"))
		return
	}

	key := param[0]

	//delete in main map table
	dbMap, err := m.MapModel.Delete(key)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	//update audit log record is latest = false
	_, err = m.AuditModel.UpdateIsLatest(key)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	//set seed for random number
	auditReq := &dto.AuditLogReq{
		Timestamp:      int(time.Now().Unix()),
		Map_id:         dbMap.ID,
		Key:            dbMap.Key,
		Original_value: dbMap.Value,
		New_value:      0,
		Action:         "deleted",
		Is_latest:      true,
		User_id:        rand.Intn(100-1) + 1,
	}

	//insert delete operation into audit log table
	_, err = m.AuditModel.Create(auditReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessRes(w, "Deleted successfully", dbMap)
}

func (m *MapHandler) Revert(w http.ResponseWriter, r *http.Request) {
	revReq := &dto.RevertReq{}
	err := json.NewDecoder(r.Body).Decode(revReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusBadRequest, err)
		return
	}

	//get record from audit log table by timpstamp
	foundLog, err := m.AuditModel.GetByTimstamp(revReq.Timestamp)
	if err != nil {
		helper.ErrorRes(w, http.StatusBadRequest, err)
		return
	}

	//get ori data from map table
	dbOriMap, err := m.MapModel.GetByKey(foundLog.Key)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	//update in the main map table
	mapReq := &dto.MapReq{
		Key:   foundLog.Key,
		Value: foundLog.New_value,
	}
	dbUpdatedMap, err := m.MapModel.Update(mapReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	//update audit log record is latest = false
	_, err = m.AuditModel.UpdateIsLatest(mapReq.Key)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	//insert in audit log table
	auditReq := &dto.AuditLogReq{
		Timestamp:      int(time.Now().Unix()),
		Map_id:         dbUpdatedMap.ID,
		Key:            dbUpdatedMap.Key,
		Original_value: dbOriMap.Value,
		New_value:      dbUpdatedMap.Value,
		Action:         "reverted",
		Is_latest:      true,
		User_id:        generateRandom(),
	}

	//insert update operation into audit log table
	_, err = m.AuditModel.Create(auditReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessRes(w, "Reverted successfully", dbUpdatedMap)
}

func generateRandom() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100-1) + 1
}
