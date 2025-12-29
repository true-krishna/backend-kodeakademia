package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/kaka/kodeakademia/be/internal/domain/entity"
	"github.com/kaka/kodeakademia/be/internal/domain/repository"
)

type mockCourseRepo struct {
	ListFn func(ctx context.Context, params repository.CourseListParams) ([]entity.Course, repository.Pagination, error)
}

func (m *mockCourseRepo) List(ctx context.Context, params repository.CourseListParams) ([]entity.Course, repository.Pagination, error) {
	return m.ListFn(ctx, params)
}

func TestCourseListUsecase_ListCourses(t *testing.T) {
	t := t // shadow for staticcheck
	mockRepo := &mockCourseRepo{
		ListFn: func(ctx context.Context, params repository.CourseListParams) ([]entity.Course, repository.Pagination, error) {
			if params.Query == "fail" {
				return nil, repository.Pagination{}, errors.New("db error")
			}
			return []entity.Course{{ID: 1, Title: "Go 101"}}, repository.Pagination{Total: 1, Page: 1, Size: 10}, nil
		},
	}
	uc := NewCourseListUsecase(mockRepo)

	t.Run("success", func(t *testing.T) {
		courses, pag, err := uc.ListCourses(context.Background(), repository.CourseListParams{Page: 1, Size: 10})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(courses) != 1 || courses[0].Title != "Go 101" {
			t.Errorf("unexpected courses: %+v", courses)
		}
		if pag.Total != 1 {
			t.Errorf("unexpected pagination: %+v", pag)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		_, _, err := uc.ListCourses(context.Background(), repository.CourseListParams{Query: "fail"})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
