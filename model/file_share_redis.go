package model

import "encoding/json"

type FileShareRedis struct {
	ID        int
	FileID    int
	URL       string
	OpenRate  int
	ExpiresAt int64
}

func (f FileShareRedis) MarshalBinary() (data []byte, err error) {
	bytes, err := json.Marshal(f)
	return bytes, err
}
