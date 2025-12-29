package usecase

import (
	"context"
	"github.com/kaka/kodeakademia/be/internal/domain/entity"
	"github.com/kaka/kodeakademia/be/internal/domain/repository"
)

type CourseListUsecase struct {
	CourseRepo repository.CourseRepository
}

func NewCourseListUsecase(repo repository.CourseRepository) *CourseListUsecase {
	return &CourseListUsecase{CourseRepo: repo}
}

func (u *CourseListUsecase) ListCourses(ctx context.Context, params repository.CourseListParams) ([]entity.Course, repository.Pagination, error) {
	return u.CourseRepo.List(ctx, params)
}
