package repentities

import (
	logicentities "svalka-service/internal/logic/entities"
)

type UserInfo struct {
	Name     string `db:"name"`
	Position string `db:"position"`
	Role     string `db:"role"`
	Phone    string `db:"phone"`
	Email    string `db:"email"`
}

type UserCreate struct {
	UserInfo
	Password string `db:"password"`
}

type User struct {
	ID uint64 `db:"id"`
	UserInfo
}

type UserWithPass struct {
	User
	Password string `db:"password"`
}

type UserConverter struct {
}

// UserWithPassToLogic ...
func (c UserConverter) UserWithPassToLogic(v UserWithPass) logicentities.UserWithPass {
	return logicentities.UserWithPass{
		User:     c.ToLogic(v.User),
		Password: v.Password,
	}
}

// userInfoToLogic ...
func (c UserConverter) userInfoToLogic(v UserInfo) logicentities.UserInfo {
	return logicentities.UserInfo{
		Name:     v.Name,
		Position: v.Position,
		Phone:    v.Phone,
		Role:     logicentities.UserRoleFromString(v.Role),
		Email:    v.Email,
	}
}

// ToLogic ...
func (c UserConverter) ToLogic(v User) logicentities.User {
	return logicentities.User{
		ID:       v.ID,
		UserInfo: c.userInfoToLogic(v.UserInfo),
	}
}

// ToLogic ...
func (c UserConverter) CreateToDB(v logicentities.UserCreate) UserCreate {
	return UserCreate{
		UserInfo: UserInfo{
			Name:     v.Name,
			Position: v.Position,
			Phone:    v.Phone,
			Email:    v.Email,
			Role:     v.Role.ToString(),
		},
		Password: v.Password,
	}
}
