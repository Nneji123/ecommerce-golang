package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			log.Infof("%s %s %d %s %s",
				req.Method,
				req.RequestURI,
				res.Status,
				time.Since(start).String(),
				c.RealIP(),
			)

			return err
		}
	}
}
