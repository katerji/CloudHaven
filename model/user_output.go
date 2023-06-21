package model

type UserOutput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (output UserOutput) ToUser() User {
	return User(output)
}
