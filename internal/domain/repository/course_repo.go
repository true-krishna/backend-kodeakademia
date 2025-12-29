package repository

import (
	"context"
	"github.com/kaka/kodeakademia/be/internal/domain/entity"
)

type CourseListParams struct {
	Page     int
	Size     int
	Query    string
	Category string
}

type Pagination struct {
	Total int
	Page  int
	Size  int
}

type CourseRepository interface {
	List(ctx context.Context, params CourseListParams) ([]entity.Course, Pagination, error)
}
