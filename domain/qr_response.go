package domain

type QrResponse struct {
	UserID string `json:"userId"`
	QrURL  string `json:"qrUrl"`
}
