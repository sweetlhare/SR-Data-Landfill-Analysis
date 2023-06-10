package logicentities

type Session struct {
	ID string `json:"session_id"`
}

type SessionCredentials struct {
	UserID uint64   `json:"id"`
	Role   UserRole `json:"role"`
}
