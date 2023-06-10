package sessionrep

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"svalka-service/internal/config"
	logicentities "svalka-service/internal/logic/entities"
	logicinterfaces "svalka-service/internal/logic/interfaces"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type sessionRep struct {
	client *redis.Client
	config config.RedisRepConfig
}

// NewSessionRep ...
func NewSessionRep(_ context.Context) (logicinterfaces.SessionRepository, error) {
	sessionRep := &sessionRep{}
	sessionRep.config = getRedisConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     sessionRep.config.RedisAddr(),
		Password: sessionRep.config.RedisPassword(),
		DB:       int(sessionRep.config.RedisDb()),
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(fmt.Errorf("redis error: %s", err.Error()))
	}
	sessionRep.client = client
	return sessionRep, nil
}

// CreateSession ...
func (r sessionRep) CreateSession(creds logicentities.SessionCredentials) (session *logicentities.Session, err error) {
	// if exist
	key := fmt.Sprint(creds.UserID)
	existingSessionID, _ := r.client.Get(key).Result()
	if existingSessionID != "" {
		return &logicentities.Session{
			ID: existingSessionID,
		}, nil
	}

	// else
	sessionID := generateSessionID(creds.UserID)
	value := strings.Join(
		[]string{
			key,                   // 0  user_id
			creds.Role.ToString(), // 1  user_role
		},
		"_",
	)

	// save session by user_id
	err = r.client.Set(key, sessionID, r.config.RedisSessionTimeout()).Err()
	if err != nil {
		return nil, err
	}

	// save creds by session_id
	err = r.client.Set(sessionID, value, r.config.RedisSessionTimeout()).Err()
	if err != nil {
		return nil, err
	}

	session = &logicentities.Session{
		ID: sessionID,
	}

	return session, nil
}

// ValidateSession ...
func (r sessionRep) ValidateSession(sessionID string) (bool, error) {
	exists, err := r.client.Exists(sessionID).Result()
	if err != nil {
		return false, err
	}

	if exists == 0 {
		return false, nil
	}

	return true, nil
}

func generateSessionID(userID uint64) string {
	return uuid.New().String()
}

// GetSessionCredentials ...
func (r sessionRep) GetSessionCredentials(sessionID string) (*logicentities.SessionCredentials, error) {
	var creds logicentities.SessionCredentials
	ok, err := r.ValidateSession(sessionID)
	if !ok {
		return nil, errors.New("invalid session")
	}
	if err != nil {
		return nil, err
	}

	value, err := r.client.Get(sessionID).Result()
	if err != nil {
		return nil, err
	}

	s := strings.Split(value, "_")
	userID, err := strconv.ParseUint(s[0], 10, 64) // s[0] - user_id
	if err != nil {
		return nil, err
	}

	creds.Role = logicentities.UserRoleFromString(s[1]) // s[1] - user_role
	creds.UserID = userID
	return &creds, nil
}
