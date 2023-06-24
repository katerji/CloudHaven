package service

import (
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
	return cache.GetRedisClient().Del(model.GetFileShareRedisKey(fileShareID))
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
	keys := cache.GetRedisClient().Keys(getRedisPrefix())
	fileShareMaps := cache.GetRedisClient().GetMulti(keys)

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
	fileShareMap := cache.GetRedisClient().Get(model.GetFileShareRedisKey(fileShareInput.ID))

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
	err := cache.GetRedisClient().Set(model.GetFileShareRedisKey(fileShare.ID), fileShare.ToRedisMap(), model.FileShareRedisExpiry)
	if err != nil {
		fmt.Println(err)
		return errors.New("error creating file share")
	}
	return nil
}

func (service fileShareService) IncrementOpenRate(fileShareInput model.FileShareInput) bool {
	fileShareString := cache.GetRedisClient().Get(model.GetFileShareRedisKey(fileShareInput.ID))

	fileShare := model.FileShare{}
	err := fileShare.Unmarshal([]byte(fileShareString))
	if err != nil {
		fmt.Println(err)
		return false
	}
	fileShare.OpenRate++

	ttl := cache.GetRedisClient().TTL(model.GetFileShareRedisKey(fileShareInput.ID))
	return cache.GetRedisClient().Set(model.GetFileShareRedisKey(fileShareInput.ID), fileShare.ToRedisMap(), ttl) == nil
}

func (service fileShareService) GetFileShares(fileID, userID int) ([]model.FileShare, error) {
	ownerID, err := GetFileService().GetFileOwner(fileID)
	if err != nil {
		return []model.FileShare{}, err
	}
	if ownerID != userID {
		return []model.FileShare{}, errors.New("unauthorized")
	}
	fileShares, err := service.fetchFileSharesByID(fileID)
	if err != nil {
		return []model.FileShare{}, errors.New("error fetching file shares")
	}
	fileSharesMap := make(map[int]model.FileShare)
	for _, fileShare := range fileShares {
		fileSharesMap[fileShare.ID] = fileShare
	}
	cachedFileShares := service.getFileSharesFromCache(fileShares)

	fileSharesToReturn := cachedFileShares
	for _, fileShare := range cachedFileShares {
		if _, ok := fileSharesMap[fileShare.ID]; !ok {
			fileSharesToReturn = append(fileSharesToReturn, fileShare)
		}
	}
	return fileSharesToReturn, nil
}

func (service fileShareService) getFileSharesFromCache(fileSharesDB []model.FileShare) []model.FileShare {
	var fileShares []model.FileShare
	for _, fileShareDB := range fileSharesDB {
		fileShareString := cache.GetRedisClient().Get(model.GetFileShareRedisKey(fileShareDB.ID))
		if fileShareString == "" {
			continue
		}
		fileShare := model.FileShare{}
		err := fileShare.Unmarshal([]byte(fileShareString))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fileShares = append(fileShares, fileShare)
	}
	return fileShares
}

func (service fileShareService) fetchFileSharesByID(fileID int) ([]model.FileShare, error) {
	rows, err := db.GetDbInstance().Query(query.FetchFileSharesQuery, fileID)
	if err != nil {
		fmt.Println(err)
		return []model.FileShare{}, errors.New("error fetching file shares")
	}
	defer rows.Close()

	var fileShares []model.FileShare
	for rows.Next() {
		fileShare := model.FileShare{}
		var expiresAtUnix int64
		err := rows.Scan(&fileShare.ID, &fileShare.FileID, &fileShare.URL, &expiresAtUnix, &fileShare.OpenRate)
		if err != nil {
			fmt.Println(err)
			return []model.FileShare{}, errors.New("error fetching file shares")
		}
		fileShare.ExpiresAt = time.Unix(expiresAtUnix, 0)
		fileShares = append(fileShares, fileShare)
	}
	return fileShares, nil
}

func getRedisPrefix() string {
	sharePrefix := model.FileShareRedisPrefix

	return fmt.Sprintf("%s*", sharePrefix)
}
