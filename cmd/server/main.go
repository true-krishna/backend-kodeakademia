package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/labstack/echo/v4"

    httptransport "github.com/kaka/kodeakademia/be/internal/transport/httptransport"
)

func main() {
    addr := ":8080"
    if p := os.Getenv("PORT"); p != "" {
        addr = fmt.Sprintf(":%s", p)
    }

    e := echo.New()

    // health route
    e.GET("/health", func(c echo.Context) error {
        return httptransport.Health(c)
    })

    log.Printf("starting server on %s", addr)
    if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}
