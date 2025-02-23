package models

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/mail"
	"time"
)

type MissingEntityError struct {
	msg        string
	identifier string
}

func (mee *MissingEntityError) Error() string {
	return mee.msg
}

func (mee *MissingEntityError) Identifier() string {
	return mee.identifier
}

func NewMissingEntityError(msg string, identifier string) *MissingEntityError {
	return &MissingEntityError{
		msg:        msg,
		identifier: identifier,
	}
}

type DateTime time.Time

func (op *DateTime) UnmarshalJSON(bytes []byte) error {
	var strVal string
	err := json.Unmarshal(bytes, &strVal)
	if err != nil {
		return err
	}
	dt, err := time.ParseInLocation("2006/01/02 15:04:05", strVal, time.UTC)
	if err != nil {
		slog.Warn(fmt.Sprintf("Could not parse invalid datetime %s", strVal), "Source", "DateTimeUnmarhaller")
		return err
	}
	*op = DateTime(dt)
	return nil
}

type EmailAddress struct {
	StringEmail string
	ParsedEmail mail.Address
}

func (ea *EmailAddress) UnmarshalJSON(bytes []byte) error {
	var strVal string
	err := json.Unmarshal(bytes, &strVal)
	if err != nil {
		return err
	}
	addr, err := mail.ParseAddress(strVal)
	if err != nil {
		slog.Warn(fmt.Sprintf("Could not parse invalid email %s", strVal), "Source", "EmailUnmarshaller")
		return err
	}
	*ea = EmailAddress{
		StringEmail: strVal,
		ParsedEmail: *addr,
	}
	return nil
}
