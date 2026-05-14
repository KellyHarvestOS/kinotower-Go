package models

type SignupRequest struct {
	FIO      string `json:"fio"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
	GenderID int    `json:"gender_id"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	FIO      string `json:"fio"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	GenderID int    `json:"gender_id"`
}

type CreateReviewRequest struct {
	FilmID  int    `json:"film_id"`
	Message string `json:"message"`
}

type CreateRatingRequest struct {
	FilmID int `json:"film_id"`
	Ball   int `json:"ball"`
}

type AuthResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
	ID     int    `json:"id"`
	FIO    string `json:"fio"`
}
