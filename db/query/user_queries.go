package query

const (
	InsertUserQuery     = "INSERT INTO user (email, name, password) VALUES (?, ?, ?)"
	GetUserByEmailQuery = "SELECT id, name, email, password FROM user WHERE email = ?"
	GetUsersQuery       = "SELECT id, name, email FROM user"
)
