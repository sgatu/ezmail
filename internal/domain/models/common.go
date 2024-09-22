package models

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
