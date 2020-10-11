package helper

import (
	"encoding/json"
	"net/http"

	model "github.com/gold-kou/ToeBeans/app/domain/model/http"
	"github.com/opentracing/opentracing-go/log"
)

func ResponseSimpleSuccess(w http.ResponseWriter) {
	resp := model.ResponseSimpleSuccess{
		Status:  http.StatusOK,
		Message: "success",
	}
	w.Header().Set(HeaderKeyContentType, "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error(err)
		panic(err.Error())
	}
}

func ResponseBadRequest(w http.ResponseWriter, message string, errors ...map[string]string) {
	resp := model.ResponseBadRequest{
		Status:  http.StatusBadRequest,
		Message: message,
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	w.Header().Set(HeaderKeyCacheControl, HeaderValueNoStore)
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error(err)
		panic(err.Error())
	}
}

func ResponseUnauthorized(w http.ResponseWriter, message string) {
	resp := model.ResponseUnauthorized{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	w.Header().Set(HeaderKeyCacheControl, HeaderValueNoStore)
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error(err)
		panic(err.Error())
	}
}

func ResponseInternalServerError(w http.ResponseWriter, message string) {
	resp := model.ResponseInternalServerError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	w.Header().Set(HeaderKeyCacheControl, HeaderValueNoStore)
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error(err)
		panic(err.Error())
	}
}
