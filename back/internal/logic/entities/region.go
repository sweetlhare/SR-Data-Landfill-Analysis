package logicentities

type Region struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"  validate:"required"`
}
