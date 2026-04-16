package controller

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=72"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateArticleRequest struct {
	Title string `json:"title" validate:"required"`
	Content string `json:"content"`
	Status string `json:"status"`
}

type UpdateArticleRequest struct {
	Title *string `json:"title"`
	Content *string `json:"content"`
	Status *string `json:"status"`
}