package domain

type TokenResponse struct {
	UserID      string `json:"userId"`
	AccessToken string `json:"accessToken"`
}
