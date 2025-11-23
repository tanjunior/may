package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Centralized API error codes used in JSON responses.
const (
	CodeInternalError    = "INTERNAL_ERROR"
	CodeProductNotFound  = "PRODUCT_NOT_FOUND"
	CodeProductsNotFound = "PRODUCTS_NOT_FOUND"
	CodeInvalidRequest   = "INVALID_REQUEST"
	CodeInvalidID        = "INVALID_ID"
	CodePerPageTooLarge  = "PER_PAGE_TOO_LARGE"
)

// ErrorMessages maps error codes to default human-readable messages.
var ErrorMessages = map[string]string{
	CodeInternalError:    "internal server error",
	CodeProductNotFound:  "product not found",
	CodeProductsNotFound: "products not found",
	CodeInvalidRequest:   "invalid request",
	CodeInvalidID:        "invalid product id",
	CodePerPageTooLarge:  "per_page exceeds maximum allowed",
}

// APIError represents a structured API error.
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// NewAPIError builds an APIError from a code and optional details.
func NewAPIError(code string, details interface{}) APIError {
	msg := ErrorMessages[code]
	return APIError{Code: code, Message: msg, Details: details}
}

// respondAPIError sends an APIError using the existing respondError envelope.
func respondAPIError(c *gin.Context, status int, apiErr APIError) {
	respondError(c, status, apiErr.Code, apiErr.Message, apiErr.Details)
}

// respondErrorCode is a convenience helper: pass only the code and details,
// the message will be populated from the default messages map.
func respondErrorCode(c *gin.Context, status int, code string, details interface{}) {
	respondAPIError(c, status, NewAPIError(code, details))
}

// Typed constructors that return an APIError and associated HTTP status.
func NewBadRequest(code string, details interface{}) (APIError, int) {
	return NewAPIError(code, details), http.StatusBadRequest
}

func NewNotFound(code string, details interface{}) (APIError, int) {
	return NewAPIError(code, details), http.StatusNotFound
}

func NewInternalError(code string, details interface{}) (APIError, int) {
	return NewAPIError(code, details), http.StatusInternalServerError
}

// Convenience responder wrappers using the typed constructors above.
func RespondBadRequest(c *gin.Context, code string, details interface{}) {
	apiErr, status := NewBadRequest(code, details)
	respondAPIError(c, status, apiErr)
}

func RespondNotFound(c *gin.Context, code string, details interface{}) {
	apiErr, status := NewNotFound(code, details)
	respondAPIError(c, status, apiErr)
}

func RespondInternal(c *gin.Context, code string, details interface{}) {
	apiErr, status := NewInternalError(code, details)
	respondAPIError(c, status, apiErr)
}
