package dto

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	User        UserData `json:"user"`
}

type UserData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
