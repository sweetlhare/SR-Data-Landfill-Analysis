package pgclient

import "svalka-service/pkg/custom"

const (
	DBConfigError custom.Error = "failed to get db config"
	PgClientError custom.Error = "failed to get pg client"
	PingError     custom.Error = "ping error"
)
