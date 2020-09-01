package Repositories

import (
	"testing"
)


func TestUserRepo_GetUserByName(t *testing.T) {
	repo := NewUserRepo()
	name := "name"
	id, _ := repo.CreateUser(name, "pass")
	user, _ := repo.UserByUsername(name)

	t.Log(id)
	t.Log(user.Id)

	if user.Id != id{
		t.Errorf("wrong id")
	}

	repo.DeleteUserByID(id)
}

