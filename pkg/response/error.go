package response

import "github.com/go-playground/validator/v10"

type APIError struct {
	Type    string            `json:"type"`
	Message string            `json:"message"`
	Detail  string            `json:"detail,omitempty"`
	Status  int               `json:"status"`
	Fields  map[string]string `json:"fields,omitempty"` //note for post validation error, like wrong gmail pasword and others
}

func newAPIError(status int, errType, message, detail string) *APIError {
	return &APIError{ //note kenapa kalo di golang itu kayak begini kalo bikin constructor balikkin pointer sama dereference ?
		Type:    errType,
		Message: message,
		Detail:  detail,
		Status:  status,
	}
}

func NewValidationError(err error) *APIError {
	fields := make(map[string]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			fields[fe.Field()] = fe.Tag()
		}
	}
	return &APIError{
		Type:    "validation_error",
		Message: "mismatch data sent",
		Status:  422,
		Fields:  fields,
	}
}

func NewParamValidationError(field, issue string) *APIError {
	return &APIError{
		Type:    "validation_error",
		Message: "mismatch data sent",
		Status:  400,
		Fields: map[string]string{
			field: issue,
		},
	}
}

func ErrInternal(detail string) *APIError {
	return newAPIError(500, "internal_error", "try again later!", detail)
}

func ErrUnAuthorized(detail string) *APIError {
	return newAPIError(401, "unauthorized", "you have not logged in yet, please login", detail)
}

func ErrBadRequest(detail string) *APIError {
	return newAPIError(400, "bad_request", "wrong data sent", detail)
}
func ErrConflict(detail string) *APIError {
	return newAPIError(409, "conflict", "Data is already existed", detail)
}
func ErrTooManyRequests(detail string) *APIError {
	return newAPIError(429, "too_many_requests", "too much request, try again later", detail)
}
