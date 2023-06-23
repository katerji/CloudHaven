package query

const (
	UpsertUserFilesBaseQuery = "INSERT INTO file (name, size, content_type, owner_id, created_on, updated_on) VALUES "
	FetchUserFilesQuery      = "SELECT id, name, size, content_type, owner_id, created_on, updated_on FROM file WHERE owner_id = ?"
)
