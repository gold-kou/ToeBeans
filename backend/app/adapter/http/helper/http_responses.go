package helper

import (
	"encoding/json"
	"log"
	"net/http"

	model "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
)

func ResponseSimpleSuccess(w http.ResponseWriter) {
	resp := model.ResponseSimpleSuccess{
		Status:  http.StatusOK,
		Message: "success",
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	// TODO superfluousでる件で、試しにコメントアウトしてみる
	// w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err.Error())
	}
}

func ResponseBadRequest(w http.ResponseWriter, message string) {
	resp := model.ResponseBadRequest{
		Status:  http.StatusBadRequest,
		Message: message,
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	w.Header().Set(HeaderKeyCacheControl, HeaderValueNoStore)
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
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
		panic(err.Error())
	}
}

func ResponseNotFound(w http.ResponseWriter, message string) {
	resp := model.ResponseUnauthorized{
		Status:  http.StatusNotFound,
		Message: message,
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	w.Header().Set(HeaderKeyCacheControl, HeaderValueNoStore)
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err.Error())
	}
}

func ResponseForbidden(w http.ResponseWriter, message string) {
	resp := model.ResponseForbidden{
		Status:  http.StatusForbidden,
		Message: message,
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	w.Header().Set(HeaderKeyCacheControl, HeaderValueNoStore)
	w.WriteHeader(http.StatusForbidden)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err.Error())
	}
}

func ResponseNotAllowedMethod(w http.ResponseWriter, message string, methods []string) {
	resp := model.ResponseNotAllowedMethod{
		Status:  http.StatusMethodNotAllowed,
		Message: message,
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	w.Header().Set(HeaderKeyCacheControl, HeaderValueNoStore)
	for _, m := range methods {
		w.Header().Set(HeaderKeyAllow, m)
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err.Error())
	}
}

func ResponseConflictError(w http.ResponseWriter, message string) {
	resp := model.ResponseConflict{
		Status:  http.StatusConflict,
		Message: message,
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	w.Header().Set(HeaderKeyCacheControl, HeaderValueNoStore)
	w.WriteHeader(http.StatusConflict)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err.Error())
	}
}

func ResponseInternalServerError(w http.ResponseWriter, message string) {
	// MEMO Sentryのお金ないので、リクエストとの紐付けが難しいが、標準ログを出すしかない
	log.Println("internal server error: " + message)

	resp := model.ResponseInternalServerError{
		Status:  http.StatusInternalServerError,
		Message: message,
	}
	w.Header().Set(HeaderKeyContentType, HeaderValueApplicationJSON)
	w.Header().Set(HeaderKeyCacheControl, HeaderValueNoStore)
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err.Error())
	}
}
