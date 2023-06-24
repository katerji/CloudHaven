package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/katerji/UserAuthKit/cache"
	"github.com/katerji/UserAuthKit/db"
	"github.com/katerji/UserAuthKit/db/query"
	"github.com/katerji/UserAuthKit/model"
	"time"
)

type fileShareService struct{}

func (service fileShareService) Insert(fileShareInput model.FileShareInput) (int, error) {
	insertId, err := db.GetDbInstance().Insert(query.InsertFileShareQuery, fileShareInput.FileID, fileShareInput.URL, fileShareInput.ExpiresAt.Unix())
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("error sharing file")
	}
	return insertId, nil
}

func (service fileShareService) updateOpenRate(fileShareInput model.FileShareInput) bool {
	return db.GetDbInstance().Exec(query.UpdateFileShareOpenRateQuery, fileShareInput.OpenRate, fileShareInput.ID)
}

func (service fileShareService) deleteCache(fileShareID int) bool {
	return cache.GetRedisClient().Del(context.Background(), model.GetFileShareRedisKey(fileShareID)).Err() == nil
}

func (service fileShareService) SyncOpenRates() {
	fileSharesToSync := service.getFileSharesToSync()
	for _, fileShare := range fileSharesToSync {
		fileShareInput := model.FileShareInput{
			ID:       fileShare.ID,
			OpenRate: fileShare.OpenRate,
		}
		go func() {
			ok := service.updateOpenRate(fileShareInput)
			if ok {
				service.deleteCache(fileShareInput.ID)
			}
		}()
	}
}

func (service fileShareService) getFileSharesToSync() []model.FileShare {
	keys := cache.GetRedisClient().Keys(context.Background(), getRedisPrefix()).Val()
	fileShareMaps := cache.GetRedisClient().MGet(context.Background(), keys...).Val()

	var fileSharesToSync []model.FileShare
	for _, fileShareMap := range fileShareMaps {
		if fileShareMap == nil {
			continue
		}
		fileShare := model.FileShare{}
		err := fileShare.Unmarshal([]byte(fileShareMap.(string)))
		if err != nil {
			fmt.Println(err)
			continue
		}
		didURLExpire := time.Now().After(fileShare.ExpiresAt)
		if didURLExpire {
			fileSharesToSync = append(fileSharesToSync, fileShare)
		}
	}

	return fileSharesToSync
}

func (service fileShareService) GetURL(fileShareInput model.FileShareInput) (string, error) {
	fileShareMap := cache.GetRedisClient().Get(context.Background(), model.GetFileShareRedisKey(fileShareInput.ID)).Val()

	fileShare := &model.FileShare{}
	err := fileShare.Unmarshal([]byte(fileShareMap))
	if err != nil {
		fmt.Println(err)
		return "", errors.New("file share not found")
	}
	if time.Now().After(fileShare.ExpiresAt) {
		return "", errors.New("url expired")
	}

	return fileShare.URL, nil
}

func (service fileShareService) SetCache(fileShare model.FileShare) error {
	err := cache.GetRedisClient().Set(context.Background(), model.GetFileShareRedisKey(fileShare.ID), fileShare.ToRedisMap(), model.FileShareRedisExpiry).Err()
	if err != nil {
		fmt.Println(err)
		return errors.New("error create file share")
	}
	return nil
}

func (service fileShareService) IncrementOpenRate(fileShareInput model.FileShareInput) {
	fileShareMap := cache.GetRedisClient().Get(context.Background(), model.GetFileShareRedisKey(fileShareInput.ID)).Val()

	fileShare := model.FileShare{}
	err := fileShare.Unmarshal([]byte(fileShareMap))
	if err != nil {
		fmt.Println(err)
		return
	}
	fileShare.OpenRate++

	ttl := cache.GetRedisClient().TTL(context.Background(), model.GetFileShareRedisKey(fileShareInput.ID)).Val()
	cache.GetRedisClient().Set(context.Background(), model.GetFileShareRedisKey(fileShareInput.ID), fileShare.ToRedisMap(), ttl)
}

func getRedisPrefix() string {
	sharePrefix := model.FileShareRedisPrefix

	return fmt.Sprintf("%s*", sharePrefix)
}
