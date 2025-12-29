package dto

type CourseListRequest struct {
	Page     int    `query:"page" validate:"min=1"`
	Size     int    `query:"size" validate:"min=1,max=100"`
	Query    string `query:"q"`
	Category string `query:"category"`
}

type CourseListItem struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	OwnerID  int64  `json:"owner_id"`
	Price    int64  `json:"price"`
	CreatedAt string `json:"created_at"`
}

type CourseListResponse struct {
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
	Items []CourseListItem  `json:"items"`
}
