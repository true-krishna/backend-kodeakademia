package http

import (
	"database/sql"
	"github.com/kaka/kodeakademia/be/internal/adapters/postgres"
	"github.com/kaka/kodeakademia/be/internal/transport/http/handlers"
	"github.com/kaka/kodeakademia/be/internal/usecase"

	"github.com/labstack/echo/v4"
)

func RegisterCourseRoutes(e *echo.Echo, db *sql.DB) {
	courseRepo := postgres.NewPostgresCourseRepository(db)
	uc := usecase.NewCourseListUsecase(courseRepo)
	h := handlers.NewCourseHandler(uc)
	e.GET("/courses", h.ListCourses)
}
