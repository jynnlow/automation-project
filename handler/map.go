package handler

import (
	"automation-project/dto"
	"automation-project/helper"
	"automation-project/model"
	"encoding/json"
	"errors"
	"net/http"
)

type MapHandlerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type MapHandler struct {
	MapModel *model.MapCRUD
}

func (m *MapHandler) Create(w http.ResponseWriter, r *http.Request) {
	mapReq := &dto.MapReq{}
	err := json.NewDecoder(r.Body).Decode(mapReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusBadRequest, err)
		return
	}

	dbMap, err := m.MapModel.Create(mapReq)
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

	helper.SuccessRes(w, "Created successfully", dbMapList)
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

	dbMap, err := m.MapModel.Update(mapReq)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessRes(w, "Updated Successfully", dbMap)
}

func (m *MapHandler) Delete(w http.ResponseWriter, r *http.Request) {
	param, exist := r.URL.Query()["key"]
	if !exist || len(param[0]) < 1 {
		helper.ErrorRes(w, http.StatusBadRequest, errors.New("param not exist"))
		return
	}

	key := param[0]

	dbMap, err := m.MapModel.Delete(key)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessRes(w, "Key found", dbMap)
}
