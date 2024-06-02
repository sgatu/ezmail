package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type BaseError struct {
	Context       map[string]string `json:"context"`
	Message       string            `json:"message"`
	ErrIdentifier string            `json:"code"`
	Code          int               `json:"-"`
}

type entityCreated struct {
	Id         string `json:"id"`
	EntityType string `json:"type"`
}

func (be BaseError) Error() string {
	return be.Message
}

func InvalidRequestBodyError() BaseError {
	return BaseError{
		Context:       make(map[string]string),
		Message:       "Invalid request body",
		ErrIdentifier: "ERR_INVALID_BODY",
		Code:          http.StatusBadRequest,
	}
}

func EntityNotFoundError(entityType string) BaseError {
	return BaseError{
		Context:       make(map[string]string),
		Message:       fmt.Sprintf("%s not found", cases.Title(language.Und, cases.NoLower).String(entityType)),
		ErrIdentifier: fmt.Sprintf("ERR_NOT_FOUND_%s", strings.ToUpper(entityType)),
		Code:          http.StatusNotFound,
	}
}

func UnauthorizedError() BaseError {
	return BaseError{
		Context:       make(map[string]string),
		Message:       "Not authorized",
		ErrIdentifier: "ERR_NOT_AUTHORIZED",
		Code:          http.StatusUnauthorized,
	}
}

func InternalServerError(err error) BaseError {
	return BaseError{
		Context:       make(map[string]string),
		Message:       fmt.Sprintf("Internal server error. Info: %s", err.Error()),
		ErrIdentifier: "ERR_GENERIC",
		Code:          http.StatusInternalServerError,
	}
}

func rawMessage(message []byte, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Write(message)
}

func ReturnReponse(response any, statusCode int, w http.ResponseWriter) {
	jsonData, err := json.Marshal(response)
	if err != nil {
		rawMessage([]byte("Failed to serialize error"), http.StatusInternalServerError, w)
		return
	}
	rawMessage(jsonData, statusCode, w)
}

func OkResponse(response any, w http.ResponseWriter) {
	ReturnReponse(response, http.StatusOK, w)
}

func CreatedResponse(response any, w http.ResponseWriter) {
	ReturnReponse(response, http.StatusCreated, w)
}

func ErrorResponse(theError error, w http.ResponseWriter) {
	switch e := theError.(type) {
	case BaseError:
		ReturnReponse(e, e.Code, w)
	default:
		ReturnReponse(BaseError{
			Context:       make(map[string]string),
			Message:       e.Error(),
			ErrIdentifier: "ERR_GENERIC",
			Code:          http.StatusInternalServerError,
		}, http.StatusInternalServerError, w)
	}
}

func EntityCreated(id string, entityType string, w http.ResponseWriter) {
	msg := entityCreated{
		Id:         id,
		EntityType: entityType,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		rawMessage([]byte("Could not serialize data"), http.StatusInternalServerError, w)
	}
	rawMessage(jsonData, http.StatusCreated, w)
}
