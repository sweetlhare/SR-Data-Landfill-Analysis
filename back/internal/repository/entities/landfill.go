package repentities

import (
	"database/sql"
)

type Landfill struct {
	ID               int64          `db:"id"`
	City             sql.NullString `db:"city"`
	RegionID         int64          `db:"region_id"`
	Illegal          bool           `db:"illegal"`
	Address          sql.NullString `db:"address"`
	Coordinates      string         `db:"coordinates"`
	PreviewImagePath sql.NullString `db:"preview_image_path"`
	Manager          sql.NullString `db:"manager"`
	Type             sql.NullString `db:"type"`
}

type LandfillCreate struct {
	City             sql.NullString `db:"city"`
	RegionID         int64          `db:"region_id"`
	Illegal          bool           `db:"illegal"`
	Address          sql.NullString `db:"address"`
	Coordinates      string         `db:"coordinates"`
	PreviewImagePath sql.NullString `db:"preview_image_path"`
	Manager          sql.NullString `db:"manager"`
	Type             sql.NullString `db:"type"`
}
