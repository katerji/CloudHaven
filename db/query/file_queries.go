package query

const (
	UpsertUserFilesBaseQuery     = "INSERT INTO file (name, size, content_type, owner_id, created_on, updated_on) VALUES "
	DeleteUserFilesBaseQuery     = "DELETE FROM file WHERE owner_id = ? AND name IN (%s)"
	FetchUserFilesQuery          = "SELECT id, name, size, content_type, owner_id, created_on, updated_on FROM file WHERE owner_id = ?"
	FetchFileQuery               = "SELECT id, name, size, content_type, owner_id, created_on, updated_on FROM file WHERE owner_id = ? AND name = ?"
	FetchFileOwnerQuery          = "SELECT owner_id FROM file WHERE id = ?"
	FetchFileSharesQuery         = "SELECT id, file_id, url, expires_at, open_rate FROM file_share WHERE file_id = ?"
	InsertFileShareQuery         = "INSERT INTO file_share (file_id, url, expires_at) VALUES (?, ?, ?)"
	UpdateFileShareOpenRateQuery = "UPDATE file_share SET open_rate = ? WHERE id = ?"
)
