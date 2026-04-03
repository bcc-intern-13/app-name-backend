package response

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Type    string            `json:"type"`
	Message string            `json:"message"`
	Detail  string            `json:"detail,omitempty"`
	Status  int               `json:"status"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func newAPIError(status int, errType, message, detail string) *APIError {
	return &APIError{
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
			fields[fe.Field()] = friendlyMessage(fe)
		}
	}
	return &APIError{
		Type:    "validation_error",
		Message: "Some fields are invalid or missing",
		Status:  422,
		Fields:  fields,
	}
}

func friendlyMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field()) //to lower case so can give consistent variable to frontend
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return field + " must be at least " + fe.Param() + " characters"
	case "max":
		return field + " must be at most " + fe.Param() + " characters"
	default:
		return field + " is invalid"
	}
}

func NewParamValidationError(field, issue string) *APIError {
	return &APIError{
		Type:    "validation_error",
		Message: "Some fields are invalid or missing",
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
	return newAPIError(400, "bad_request", "The data you sent is incorrect, please check again", detail)
}
func ErrConflict(detail string) *APIError {
	return newAPIError(409, "conflict", "Data is already existed", detail)
}
func ErrTooManyRequests(detail string) *APIError {
	return newAPIError(429, "too_many_requests", "too much request, try again later", detail)
}

func ErrNotFound(detail string) *APIError {
	return newAPIError(404, "not_found", "data not found", detail)
}
func ErrForbidden(detail string) *APIError {
	return newAPIError(403, "forbidden", "you are not allowed to access this resource", detail)
}
