package repentities

import logicentities "svalka-service/internal/logic/entities"

type Violation struct {
	ID            uint64 `db:"id"`
	Title         string `db:"description"`
	DefaultStatus bool   `db:"default_status"`
	Status        bool   `db:"status"`
}

type ViolationResponseAll struct {
	ID            uint64 `db:"id"`
	Title         string `db:"description"`
	DefaultStatus bool   `db:"default_status"`
}

type ViolationCreate struct {
	Description   string `db:"description"`
	DefaultStatus bool   `db:"default_status"`
}

type ViolationConverter struct {
}

// ToLogic ...
func (c ViolationConverter) ToLogic(v Violation) logicentities.Violation {
	return logicentities.Violation(v)
}

// ToLogicArray ...
func (c ViolationConverter) ToLogicArray(vs ...Violation) (r []logicentities.Violation) {
	for _, v := range vs {
		r = append(r, c.ToLogic(v))
	}
	return r
}
