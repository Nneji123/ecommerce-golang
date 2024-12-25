package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			code    = http.StatusInternalServerError
			message string
		)

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		} else {
			message = err.Error()
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				c.NoContent(code)
			} else {
				c.JSON(code, ErrorResponse{
					Message: message,
					Code:    code,
				})
			}
		}
	}
}