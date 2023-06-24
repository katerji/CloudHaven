package model

import "time"

type FileShareInput struct {
	ID        int
	FileID    int
	URL       string
	ExpiresAt time.Time
	OpenRate  int
}
