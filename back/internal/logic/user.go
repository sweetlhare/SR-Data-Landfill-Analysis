package logic

import (
	"context"
	"errors"
	logicEntities "svalka-service/internal/logic/entities"
	logicentities "svalka-service/internal/logic/entities"
)

// CreateUser ...
func (l logic) CreateUser(ctx context.Context, user logicEntities.UserCreate) error {
	err := l.validator.CommonValidation(user)
	if err != nil {
		return err
	}
	err = l.rep.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// Login ...
func (l logic) Login(ctx context.Context, login logicEntities.Login) (*logicEntities.Session, error) {
	err := l.validator.CommonValidation(login)
	if err != nil {
		return nil, err
	}

	user, err := l.rep.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("Login error: user does not exist.")
	}

	if user.Password != login.Password {
		return nil, errors.New("Login error: Invalid credentials.")
	}

	session, err := l.sessRep.CreateSession(logicEntities.SessionCredentials{
		UserID: user.ID,
		Role:   user.Role,
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (l logic) Auth(sessionID string, validRoles ...logicentities.UserRole) (*logicentities.SessionCredentials, error) {
	creds, err := l.sessRep.GetSessionCredentials(sessionID)
	if err != nil {
		return nil, err
	}

	if creds == nil {
		return nil, errors.New("Auth error: session creds error")
	}

	if creds.Role.IsValid() && creds.Role.Equel(validRoles...) {
		return creds, nil
	}

	return nil, errors.New("Auth error: you don't have sufficient permissions")
}
