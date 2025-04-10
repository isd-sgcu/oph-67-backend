package domain

type CertTokenResponse struct {
	UserID    string `json:"userId"`
	CertToken string `json:"certToken"`
}
