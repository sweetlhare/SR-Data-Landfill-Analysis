LOCAL_MIGRATION_DIR=./migrations
LOCAL_MIGRATION_DSN="host=95.163.237.107 port=5432 dbname=db user=user password=rL1jD3xW8p sslmode=disable"
#host=localhost port=54322 dbname=svalka user=svalka-user password=svalka-password sslmode=disable
#host=95.163.237.107 port=5432 dbname=db user=user password=rL1jD3xW8p sslmode=disable
local-migration-status:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v