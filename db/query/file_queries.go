package query

const (
	UpsertUserFilesBaseQuery = "INSERT INTO file (name, size, content_type, owner_id, created_on, updated_on) VALUES "
	DeleteUserFilesBaseQuery = "DELETE FROM file WHERE owner_id = ? AND name IN (%s)"
	FetchUserFilesQuery      = "SELECT id, name, size, content_type, owner_id, created_on, updated_on FROM file WHERE owner_id = ?"
)
