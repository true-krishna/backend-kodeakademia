package httptransport

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
    // Red: test expresses desired behavior
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // call handler
    err := Health(c)
    require.NoError(t, err)

    require.Equal(t, http.StatusOK, rec.Code)
    require.Contains(t, rec.Body.String(), "status")
}
