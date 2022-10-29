package handler

import (
	"automation-project/helper"
	"automation-project/model"
	"errors"
	"net/http"
)

type AuditHandlerInterface interface {
	List(w http.ResponseWriter, r *http.Request)
	GetListByKey(w http.ResponseWriter, r *http.Request)
}

type AuditHandler struct {
	AuditModel *model.AuditLogCRUD
}

func (a *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	dbAuditList, err := a.AuditModel.List()
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessRes(w, "Retrieved successfully", dbAuditList)
}

func (a *AuditHandler) GetListByKey(w http.ResponseWriter, r *http.Request) {
	param, exist := r.URL.Query()["key"]
	if !exist || len(param[0]) < 1 {
		helper.ErrorRes(w, http.StatusBadRequest, errors.New("param not exist"))
		return
	}

	key := param[0]
	dbAuditList, err := a.AuditModel.GetListByKey(key)
	if err != nil {
		helper.ErrorRes(w, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessRes(w, "Key found", dbAuditList)
}
