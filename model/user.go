package model

type User struct {
	ID    int
	Name  string
	Email string
}

func (user User) ToOutput() UserOutput {
	return UserOutput(user)
}
