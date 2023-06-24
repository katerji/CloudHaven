package service

import (
	"github.com/katerji/UserAuthKit/db"
	"github.com/katerji/UserAuthKit/db/query"
	"github.com/katerji/UserAuthKit/model"
)

type userService struct {
	users map[int]model.User
}

func initUserService() *userService {
	users := fetchUsers()
	return &userService{
		users: users,
	}
}

func (service *userService) GetUsers() map[int]model.User {
	return service.users
}
func (service *userService) setUsers(users map[int]model.User) {
	service.users = users
}

func fetchUsers() map[int]model.User {
	users := make(map[int]model.User)
	rows, err := db.GetDbInstance().Query(query.GetUsersQuery)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			panic(err)
		}
		users[user.ID] = user
	}
	return users
}

func (service *userService) SyncUsers() {
	users := fetchUsers()
	service.users = users
}
