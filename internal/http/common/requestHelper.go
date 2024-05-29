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

func rawMessage(message []byte, status int, w http.ResponseWriter) {
	w.WriteHeader(status)
	w.Write(message)
}

func ReturnError(err error, w http.ResponseWriter) {
	switch e := err.(type) {
	case BaseError:
		jsonData, err := json.Marshal(e)
		if err != nil {
			rawMessage([]byte("Failed to serialize error"), http.StatusInternalServerError, w)
			return
		}
		rawMessage(jsonData, e.Code, w)
	default:
		jsonData, err := json.Marshal(BaseError{
			Context:       make(map[string]string),
			Message:       e.Error(),
			ErrIdentifier: "ERR_GENERIC",
			Code:          http.StatusInternalServerError,
		})
		if err != nil {
			rawMessage([]byte("Failed to serialize error"), http.StatusInternalServerError, w)
			return
		}
		rawMessage(jsonData, http.StatusInternalServerError, w)

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
