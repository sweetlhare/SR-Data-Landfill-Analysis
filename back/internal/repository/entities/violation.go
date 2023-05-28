package repentities

type Violation struct {
	ID            int64  `db:"id"`
	Description   string `db:"description"`
	DefaultStatus bool   `db:"default_status"`
}

type ViolationCreate struct {
	Description   string `db:"description"`
	DefaultStatus bool   `db:"default_status"`
}
