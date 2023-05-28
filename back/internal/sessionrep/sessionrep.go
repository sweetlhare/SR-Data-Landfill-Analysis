package sessionrep

import (
	"context"
	"errors"
	"fmt"
	"log"
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
func (r sessionRep) CreateSession(email string) (session *logicentities.Session, err error) {

	sID := generateSessionID(email)
	session = &logicentities.Session{
		ID: sID,
	}
	key := fmt.Sprintf("session:%s", sID)
	value := email

	existingSessionID, err := r.client.Get(key).Result()
	if err == nil {
		return &logicentities.Session{
			ID: existingSessionID,
		}, nil
	}

	err = r.client.Set(key, value, r.config.RedisSessionTimeout()).Err()
	if err != nil {
		return nil, err
	}

	return session, nil
}

// ValidateSession ...
func (r sessionRep) ValidateSession(sessionID string) (bool, error) {
	key := fmt.Sprintf("session:%s", sessionID)

	exists, err := r.client.Exists(key).Result()
	if err != nil {
		return false, err
	}

	if exists == 0 {
		return false, nil
	}

	return true, nil
}

func generateSessionID(email string) string {
	return fmt.Sprintf("%s-%s", email, uuid.New().String())
}

// GetValue ...
func (r sessionRep) GetValue(sessionID string) (string, error) {
	ok, err := r.ValidateSession(sessionID)
	if !ok {
		return "", errors.New("invalid session")
	}
	if err != nil {
		return "", err
	}

	s := strings.Split(sessionID, "-")

	return s[0], nil
}
