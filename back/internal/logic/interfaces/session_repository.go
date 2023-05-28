package logicinterfaces

import logicentities "svalka-service/internal/logic/entities"

type SessionRepository interface {
	CreateSession(email string) (*logicentities.Session, error)
	ValidateSession(sessionID string) (bool, error)
	GetValue(sessionID string) (string, error)
}
