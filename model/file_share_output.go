package model

type FileShareOutput struct {
	ID        int    `json:"id"`
	FileID    int    `json:"file_id"`
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at"`
	OpenRate  int    `json:"open_rate"`
}
