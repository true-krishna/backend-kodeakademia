package httptransport

import (
    "net/http"

    "github.com/labstack/echo/v4"
)

// Health returns a simple health-check payload.
func Health(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
