package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the unified API envelope.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// OK writes a successful unified response.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

// Created writes a 201 unified response.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Code: 0, Message: "ok", Data: data})
}

// Fail writes an error unified response with the given HTTP status and business code.
func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, Response{Code: code, Message: message, Data: nil})
}
