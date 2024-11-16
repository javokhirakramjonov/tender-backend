package request_model

import "time"

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

type CreateTenderReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Budget      float64 `json:"budget"`
}

type UpdateTenderReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Budget      float64 `json:"budget"`
	Status      string `json:"status"`
}

type ValidateTenderBelongsToUserReq struct {
	ClientID string `json:"client_id"`

}

