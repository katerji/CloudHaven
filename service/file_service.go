package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/katerji/UserAuthKit/cache"
	"github.com/katerji/UserAuthKit/db"
	"github.com/katerji/UserAuthKit/db/query"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/utils"
	"time"
)

type fileService struct{}

func (service fileService) upsertUserFiles(files []model.File) bool {
	if len(files) == 0 {
		return true
	}
	upsertQuery := query.UpsertUserFilesBaseQuery
	var upsertQueryArgs []any
	for _, file := range files {
		args := []any{
			file.Name,
			file.Size,
			file.ContentType,
			file.OwnerID,
			utils.TimeToString(file.CreatedOn),
			utils.TimeToString(file.ModifiedOn),
		}
		upsertQuery += fmt.Sprintf("(%s),", db.GetQuestionMarks(len(args)))
		upsertQueryArgs = append(upsertQueryArgs, args...)
	}
	upsertQuery = upsertQuery[:len(upsertQuery)-1]
	upsertQuery += " ON DUPLICATE KEY UPDATE size=VALUES(size), content_type=VALUES(content_type), updated_on=VALUES(updated_on)"

	return db.GetDbInstance().Exec(upsertQuery, upsertQueryArgs...)
}

func (service fileService) deleteUserFiles(files []model.File) bool {
	if len(files) == 0 {
		return true
	}
	var fileNames []any
	userID := files[0].OwnerID
	for _, file := range files {
		if file.OwnerID != userID {
			return false
		}
		fileNames = append(fileNames, file.Name)
	}
	deleteQuery := fmt.Sprintf(query.DeleteUserFilesBaseQuery, db.GetQuestionMarks(len(fileNames)))
	args := []any{
		userID,
	}
	args = append(args, fileNames...)

	return db.GetDbInstance().Exec(deleteQuery, args...)
}

func (service fileService) GetUserFiles(userID int) (map[string]model.File, error) {
	rows, err := db.GetDbInstance().Query(query.FetchUserFilesQuery, userID)
	if err != nil {
		fmt.Println(err)
		return make(map[string]model.File), errors.New("error fetching user files")
	}
	files := make(map[string]model.File)
	for rows.Next() {
		var file model.File
		var createdOn string
		var modifiedOn string
		err = rows.Scan(&file.ID, &file.Name, &file.Size, &file.ContentType, &file.OwnerID, &createdOn, &modifiedOn)
		if err != nil {
			fmt.Println(err)
			return make(map[string]model.File), errors.New("error fetching user files")
		}
		file.CreatedOn, _ = time.Parse(utils.DateTimeFormat, createdOn)
		file.ModifiedOn, _ = time.Parse(utils.DateTimeFormat, modifiedOn)
		files[file.Name] = file
	}
	return files, nil
}

func (service fileService) GetFile(fileInput model.FileInput) (model.File, error) {
	row := db.GetDbInstance().QueryRow(query.FetchFileQuery, fileInput.OwnerID, fileInput.Name)
	var file model.File
	var createdOn string
	var modifiedOn string
	err := row.Scan(&file.ID, &file.Name, &file.Size, &file.ContentType, &file.OwnerID, &createdOn, &modifiedOn)
	if err != nil {
		fmt.Println(err)
		return model.File{}, errors.New("file not found")
	}
	file.CreatedOn, _ = time.Parse(utils.DateTimeFormat, createdOn)
	file.ModifiedOn, _ = time.Parse(utils.DateTimeFormat, modifiedOn)
	return file, nil
}

func (service fileService) GetFileShareURL(fileShareInput model.FileShareInput) (string, error) {
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

func (service fileService) SetFileShareCache(fileShare model.FileShare) error {
	err := cache.GetRedisClient().Set(context.Background(), model.GetFileShareRedisKey(fileShare.ID), fileShare.ToRedisMap(), model.FileShareRedisExpiry).Err()
	if err != nil {
		fmt.Println(err)
		return errors.New("error create file share")
	}
	return nil
}

func (service fileService) IncrementOpenRate(fileShareInput model.FileShareInput) {
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

func (service fileService) SyncUserFiles(userID int) error {
	gcsFiles, ok := GetGcpService().ListUserObjects(userID)
	if !ok {
		return errors.New("error syncing user files")
	}
	dbFiles, err := service.GetUserFiles(userID)
	if err != nil {
		fmt.Println(err)
		return errors.New("error syncing user files")
	}
	var newOrUpdatedFiles []model.File
	var deletedFiles []model.File
	for _, gcsFile := range gcsFiles {
		dbFile, ok := dbFiles[gcsFile.Name]
		if !ok {
			newOrUpdatedFiles = append(newOrUpdatedFiles, gcsFile)
			continue
		}
		if gcsFile.ModifiedOn.After(dbFile.ModifiedOn) {
			newOrUpdatedFiles = append(newOrUpdatedFiles, gcsFile)
		}
	}
	for _, dbFile := range dbFiles {
		_, ok := gcsFiles[dbFile.Name]
		if !ok {
			deletedFiles = append(deletedFiles, dbFile)
		}
	}
	go service.upsertUserFiles(newOrUpdatedFiles)
	go service.deleteUserFiles(deletedFiles)
	return nil
}

func (service fileService) InsertFileShare(fileShareInput model.FileShareInput) (int, error) {
	expiresAtString := fileShareInput.ExpiresAt.Format(utils.DateTimeFormat)
	insertId, err := db.GetDbInstance().Insert(query.InsertFileShareQuery, fileShareInput.FileID, fileShareInput.URL, expiresAtString)
	if err != nil {
		fmt.Println(err)
		return 0, errors.New("error sharing file")
	}
	return insertId, nil
}

func (service fileService) updateFileShareOpenRate(fileShareInput model.FileShareInput) bool {
	return db.GetDbInstance().Exec(query.UpdateFileShareOpenRateQuery, fileShareInput.OpenRate, fileShareInput.ID)
}

func (service fileService) SyncOpenRates() {
	fileSharesToSync := service.getFileSharesToSync()
	for _, fileShare := range fileSharesToSync {
		fileShareInput := model.FileShareInput{
			ID:       fileShare.ID,
			OpenRate: fileShare.OpenRate,
		}
		service.updateFileShareOpenRate(fileShareInput)
	}
}

func (service fileService) getFileSharesToSync() []model.FileShare {
	keys := cache.GetRedisClient().Keys(context.Background(), getFileSharePrefix()).Val()
	fileShareMaps := cache.GetRedisClient().MGet(context.Background(), keys...).Val()

	var fileSharesToSync []model.FileShare
	for _, fileShareMap := range fileShareMaps {
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

func getFileSharePrefix() string {
	sharePrefix := model.FileShareRedisPrefix

	return fmt.Sprintf("%s*", sharePrefix)
}
