package apperror

import "encoding/json"

var (
	ErrNotFound = NewAppError(nil, "not found", "", "US-000003")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (err *AppError) Error() string {
	return err.Message
}

func (err *AppError) Unwrap() error {
	return err.Err
}

func (err *AppError) Marshal() []byte {
	marshal, e := json.Marshal(err)
	if e != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage, code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}
