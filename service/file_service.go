package service

import (
	"errors"
	"fmt"
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
