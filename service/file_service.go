package service

import (
	"fmt"
	"github.com/katerji/UserAuthKit/db"
	"github.com/katerji/UserAuthKit/db/query"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/utils"
	"time"
)

type fileService struct{}

func (service fileService) UpsertUserFiles(files []model.File) bool {
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

func (service fileService) GetUserFiles(userID int) []model.File {
	rows, err := db.GetDbInstance().Query(query.FetchUserFilesQuery, userID)
	if err != nil {
		fmt.Println(err)
		return []model.File{}
	}
	var files []model.File
	for rows.Next() {
		var file model.File
		var createdOn string
		var modifiedOn string
		err = rows.Scan(&file.ID, &file.Name, &file.Size, &file.ContentType, &file.OwnerID, &createdOn, &modifiedOn)
		if err != nil {
			panic(err)
		}
		file.CreatedOn, _ = time.Parse(utils.DateTimeFormat, createdOn)
		file.ModifiedOn, _ = time.Parse(utils.DateTimeFormat, modifiedOn)
		files = append(files, file)
	}
	return files
}

func (service fileService) SyncUserFiles(userID int) {
	files, ok := GetGcpService().ListUserObjects(userID)
	if !ok {
		return
	}
	service.UpsertUserFiles(files)
}
