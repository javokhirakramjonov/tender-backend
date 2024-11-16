package request_model

type CreateUserReq struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type CreateNotificationReq struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}
