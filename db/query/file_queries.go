package query

const (
	UpsertUserFilesBaseQuery     = "INSERT INTO file (name, size, content_type, owner_id, created_on, updated_on) VALUES "
	DeleteUserFilesBaseQuery     = "DELETE FROM file WHERE owner_id = ? AND name IN (%s)"
	FetchUserFilesQuery          = "SELECT id, name, size, content_type, owner_id, created_on, updated_on FROM file WHERE owner_id = ?"
	FetchFileQuery               = "SELECT id, name, size, content_type, owner_id, created_on, updated_on FROM file WHERE owner_id = ? AND name = ?"
	FetchFileShareQuery          = "SELECT id, file_id, url, expires_at, FROM file_share WHERE id = ?"
	InsertFileShareQuery         = "INSERT INTO file_share (file_id, url, expires_at) VALUES (?, ?, ?)"
	UpdateFileShareOpenRateQuery = "UPDATE file_share SET open_rate = ? WHERE id = ?"
)
