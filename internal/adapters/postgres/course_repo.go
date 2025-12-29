package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/kaka/kodeakademia/be/internal/domain/entity"
	"github.com/kaka/kodeakademia/be/internal/domain/repository"
)

type PostgresCourseRepository struct {
	DB *sql.DB
}

func NewPostgresCourseRepository(db *sql.DB) *PostgresCourseRepository {
	return &PostgresCourseRepository{DB: db}
}

func (r *PostgresCourseRepository) List(ctx context.Context, params repository.CourseListParams) ([]entity.Course, repository.Pagination, error) {
	       var (
		       courses []entity.Course
		       args    []interface{}
		       where   []string
		       paramIdx = 1
		       limit   = params.Size
		       offset  = (params.Page - 1) * params.Size
	       )

	       q := `SELECT id, title, category, owner_id, price, created_at, updated_at FROM courses`
	       if params.Query != "" {
		       where = append(where, "title ILIKE $"+strconv.Itoa(paramIdx))
		       args = append(args, "%"+params.Query+"%")
		       paramIdx++
	       }
	       if params.Category != "" {
		       where = append(where, "category = $"+strconv.Itoa(paramIdx))
		       args = append(args, params.Category)
		       paramIdx++
	       }
	       if len(where) > 0 {
		       q += " WHERE " + strings.Join(where, " AND ")
	       }
	       q += " ORDER BY created_at DESC LIMIT $"+strconv.Itoa(paramIdx)+" OFFSET $"+strconv.Itoa(paramIdx+1)
	       args = append(args, limit, offset)

	       rows, err := r.DB.QueryContext(ctx, q, args...)
	       if err != nil {
		       return nil, repository.Pagination{}, err
	       }
	       defer rows.Close()

	       for rows.Next() {
		       var c entity.Course
		       if err := rows.Scan(&c.ID, &c.Title, &c.Category, &c.OwnerID, &c.Price, &c.CreatedAt, &c.UpdatedAt); err != nil {
			       return nil, repository.Pagination{}, err
		       }
		       courses = append(courses, c)
	       }

	       // Get total count
	       totalQ := "SELECT COUNT(*) FROM courses"
	       if len(where) > 0 {
		       totalQ += " WHERE " + strings.Join(where, " AND ")
	       }
	       totalArgs := args[:len(args)-2] // exclude limit/offset
	       totalRow := r.DB.QueryRowContext(ctx, totalQ, totalArgs...)
	       var total int
	       if err := totalRow.Scan(&total); err != nil {
		       return nil, repository.Pagination{}, err
	       }

	       pagination := repository.Pagination{
		       Total: total,
		       Page:  params.Page,
		       Size:  params.Size,
	       }
	       return courses, pagination, nil
}
