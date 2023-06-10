package logicinterfaces

import logicentities "svalka-service/internal/logic/entities"

type SessionRepository interface {
	CreateSession(s logicentities.SessionCredentials) (*logicentities.Session, error)
	ValidateSession(sessionID string) (bool, error)
	GetSessionCredentials(sessionID string) (*logicentities.SessionCredentials, error)
}
