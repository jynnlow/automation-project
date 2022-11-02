package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func ErrorRes(w http.ResponseWriter, status_code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	w.Write([]byte(err.Error()))
}

func SuccessRes(w http.ResponseWriter, message string, details interface{}) {
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	res := &Response{
		Message: message,
		Details: details,
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(jsonRes)
}
