package repentities

type Region struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type RegionCreate struct {
	Name string `db:"name"`
}
