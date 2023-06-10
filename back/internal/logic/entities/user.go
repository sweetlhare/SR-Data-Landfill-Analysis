package logicentities

type UserInfo struct {
	Name     string   `json:"name"`
	Position string   `json:"position"`
	Role     UserRole `json:"role"`
	Phone    string   `json:"phone"`
	Email    string   `json:"email" validate:"required"`
}

type User struct {
	ID uint64 `json:"id"`
	UserInfo
}

type UserWithPass struct {
	User
	Password string `json:"password" validate:"required"`
}

type UserCreate struct {
	UserInfo
	Password string `json:"password" validate:"required"`
}

type Login struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}
