package examples

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// EchoHandler creates http.Handler using echo framework.
//
// Routes:
//  GET /login             authenticate user and return JWT token
//  GET /restricted/hello  return "hello, world!" (requires authentication)
func EchoHandler() http.Handler {
	e := echo.New()

	e.GET("/hello", func(ctx echo.Context) error {
		name := ctx.Request().Header.Get("X-UserName")
		if name == "" {
			name = "world"
		}
		return ctx.String(http.StatusOK, fmt.Sprintf("hello, %s!", name))
	})

	return e
}
