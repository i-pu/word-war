package entity

type Result struct {
	UserID string `json:"userID" validate:"required"`
	Score  int64  `json:"score" validate:"required"`
}
