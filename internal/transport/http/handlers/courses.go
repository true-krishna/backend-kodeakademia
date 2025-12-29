package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/kaka/kodeakademia/be/internal/domain/repository"
	"github.com/kaka/kodeakademia/be/internal/transport/http/dto"
	"github.com/kaka/kodeakademia/be/internal/usecase"
)

type CourseHandler struct {
	UC *usecase.CourseListUsecase
}

func NewCourseHandler(uc *usecase.CourseListUsecase) *CourseHandler {
	return &CourseHandler{UC: uc}
}

func (h *CourseHandler) ListCourses(c echo.Context) error {
	var req dto.CourseListRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.Page == 0 {
		if p := c.QueryParam("page"); p != "" {
			if v, err := strconv.Atoi(p); err == nil {
				req.Page = v
			}
		}
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		if s := c.QueryParam("size"); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				req.Size = v
			}
		}
	}
	if req.Size == 0 {
		req.Size = 10
	}

	params := repository.CourseListParams{
		Page:     req.Page,
		Size:     req.Size,
		Query:    req.Query,
		Category: req.Category,
	}
	courses, pag, err := h.UC.ListCourses(c.Request().Context(), params)
	if err != nil {
		fmt.Printf("ListCourses error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list courses"})
	}
	items := make([]dto.CourseListItem, len(courses))
	for i, course := range courses {
		items[i] = dto.CourseListItem{
			ID:        course.ID,
			Title:     course.Title,
			Category:  course.Category,
			OwnerID:   course.OwnerID,
			Price:     course.Price,
			CreatedAt: course.CreatedAt,
		}
	}
	resp := dto.CourseListResponse{
		Total: pag.Total,
		Page:  pag.Page,
		Size:  pag.Size,
		Items: items,
	}
	return c.JSON(http.StatusOK, resp)
}
