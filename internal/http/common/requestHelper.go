package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/sgatu/ezmail/internal/domain/models"
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

type resultMessage struct {
	Message string `json:"message"`
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

func InvalidRequest() BaseError {
	return BaseError{
		Context:       make(map[string]string),
		Message:       "Invalid request",
		ErrIdentifier: "ERR_INVALID_REQUEST",
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

func OkOperation(w http.ResponseWriter) {
	OkResponse(resultMessage{Message: "Operation successful"}, w)
}

func CreatedResponse(response any, w http.ResponseWriter) {
	ReturnReponse(response, http.StatusCreated, w)
}

func ErrorResponse(theError error, w http.ResponseWriter) {
	switch e := theError.(type) {
	case BaseError:
		ReturnReponse(e, e.Code, w)
	case *models.MissingEntityError:
		ReturnReponse(BaseError{
			Context:       map[string]string{"identifier": e.Identifier()},
			Message:       e.Error(),
			ErrIdentifier: "ERR_NOT_FOUND",
			Code:          http.StatusNotFound,
		}, http.StatusNotFound, w)
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

type RegistrationMethod func(path string, handlerfunc http.HandlerFunc)

func RegisterEndpoint(
	method RegistrationMethod,
	path string,
	handlerfunc http.HandlerFunc,
	description string,
) {
	methodName := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
	methodNameParts := strings.Split(methodName, ".")
	methodName = methodNameParts[len(methodNameParts)-1]
	methodNameParts = strings.Split(methodName, "-")
	methodName = strings.ToUpper(methodNameParts[0])
	fmt.Printf("%s %s -> %s\n", methodName, path, description)
	method(path, handlerfunc)
}
