package repentities

import (
	logicentities "svalka-service/internal/logic/entities"
)

type Region struct {
	ID   *uint64 `db:"id"`
	Name string  `db:"name"`
}

type RegionCreate struct {
	Name string `db:"name"`
}

type RegionConverter struct {
}

// ToLogic ...
func (c RegionConverter) CreateToDB(a logicentities.Region) RegionCreate {
	return RegionCreate{
		Name: a.Title,
	}
}

// ToLogic ...
func (c RegionConverter) ToLogic(a *Region) logicentities.Region {
	return logicentities.Region{
		ID:    *a.ID,
		Title: a.Name,
	}
}

// ToLogic ...
func (c RegionConverter) ToLogicArray(vs ...*Region) (r []logicentities.Region) {
	for _, v := range vs {
		if v != nil {
			r = append(r, c.ToLogic(v))
		}
	}
	return r
}
