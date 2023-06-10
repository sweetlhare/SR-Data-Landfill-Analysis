package repentities

import (
	logicentities "svalka-service/internal/logic/entities"
	"svalka-service/pkg/pg"
)

type Manager struct {
	Name     string `db:"name"`
	Position string `db:"position"`
	Phone    string `db:"phone"`
	Email    string `db:"email"`
}

type ManagerConverter struct {
}

// ToLogic ...
func (c ManagerConverter) ToLogic(m Manager) logicentities.Manager {
	return logicentities.Manager(m)
}

// ToDB ...
func (c ManagerConverter) ToDB(m logicentities.Manager) string {
	return pg.StructToDBString(m)
}
