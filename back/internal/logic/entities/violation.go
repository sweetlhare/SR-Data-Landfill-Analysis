package logicentities

type Violation struct {
	ID            uint64 `json:"id"`
	Title         string `json:"title"`
	DefaultStatus bool   `json:"default"`
	Status        bool   `json:"status"`
}
