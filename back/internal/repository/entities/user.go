package repentities

import (
	"database/sql"
)

type User struct {
	ID       int64          `db:"id"`
	Name     string         `db:"name"`
	Position sql.NullString `db:"position"`
	Phone    sql.NullString `db:"phone"`
	Email    string         `db:"email"`
	Password string         `db:"password"`
}

type UserCreate struct {
	Name     string         `db:"name"`
	Position sql.NullString `db:"position"`
	Phone    sql.NullString `db:"phone"`
	Email    string         `db:"email"`
	Password string         `db:"password"`
}
