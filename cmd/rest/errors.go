package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type httpErrors struct {
	logger *slog.Logger
}

func newHTTPErrors(logger *slog.Logger) httpErrors {
	return httpErrors{
		logger: logger,
	}
}

func (h httpErrors) badRequestResponse(c *gin.Context, err error) {
	h.errorResponse(c, http.StatusBadRequest, err.Error())
}

func (h httpErrors) failedValidationResponse(c *gin.Context, errors map[string]string) {
	h.errorResponse(c, http.StatusUnprocessableEntity, errors)
}

func (h httpErrors) errorResponse(c *gin.Context, status int, message any) {
	c.IndentedJSON(status, gin.H{
		"error": message,
	})
}

func (h httpErrors) serverErrorResponse(c *gin.Context, err error) {
	h.logError(c, err)

	message := "the server encountered a problem and could not process your request"
	h.errorResponse(c, http.StatusInternalServerError, message)
}

func (h httpErrors) logError(c *gin.Context, err error) {
	var (
		method = c.Request.Method
		uri    = c.Request.URL.RequestURI()
	)

	h.logger.Error(err.Error(), "method", method, "uri", uri)
}
