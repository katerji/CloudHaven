package model

import (
	"encoding/json"
	"github.com/katerji/UserAuthKit/utils"
	"strconv"
	"time"
)

const (
	FileShareRedisKeyID        = "ID"
	FileShareRedisKeyFileID    = "FileID"
	FileShareRedisKeyURL       = "URL"
	FileShareRedisKeyExpiresAt = "ExpiresAt"
	FileShareRedisKeyOpenRate  = "OpenRate"
	FileShareRedisExpiry       = 48 * time.Hour
	FileShareRedisPrefix       = "file_share_"
)

type FileShare struct {
	ID        int
	FileID    int
	URL       string
	OpenRate  int
	ExpiresAt time.Time
}

func (f *FileShare) ToRedisMap() FileShareRedis {
	return FileShareRedis{
		ID:        f.ID,
		FileID:    f.FileID,
		URL:       f.URL,
		OpenRate:  f.OpenRate,
		ExpiresAt: f.ExpiresAt.Unix(),
	}
}

func GetFileShareRedisKey(ShareID int) string {
	return FileShareRedisPrefix + strconv.Itoa(ShareID)
}

func (f *FileShare) Unmarshal(data []byte) error {
	var fields map[string]any
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	f.ID = int(fields[FileShareRedisKeyID].(float64))
	f.FileID = int(fields[FileShareRedisKeyFileID].(float64))
	f.URL = fields[FileShareRedisKeyURL].(string)
	f.OpenRate = int(fields[FileShareRedisKeyOpenRate].(float64))

	expiresAtInt := int64(int(fields[FileShareRedisKeyExpiresAt].(float64)))
	f.ExpiresAt = time.Unix(expiresAtInt, 0)
	return nil
}

func (f *FileShare) ToOutput() FileShareOutput {
	return FileShareOutput{
		ID:        f.ID,
		FileID:    f.FileID,
		URL:       f.URL,
		ExpiresAt: f.ExpiresAt.Format(utils.DateTimeFormat),
		OpenRate:  f.OpenRate,
	}
}
