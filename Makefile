dsn?=
dbUser?=
dbPassword?=
dbName?=
dbHost?=
dbPort?=
password?=
migration_name?=
version?=


.PHONY: migrateup
create_migrate: 
	migrate create -ext sql -dir internal/db/migration -seq $(migration_name)

.PHONY: migrateup
migrateup:
	migrate -path internal/db/migration -database $(dsn) -verbose up

# dsn?= the database to run the migration. This is option is to be set from the terminal for security reasons.
# (Optional) version?=.... to rollback to a previous version in a case where a migration fails.
.PHONY: migratedown
migratedown:
	if [ $(version) ]; then \
		migrate -path internal/db/migration -database "postgresql://$(dbUser):$(dbPassword)@$(dbHost):$(dbPassword)/$(dbName)?sslmode=disable" -verbose force $(version); \
	else \
		migrate -path internal/db/migration -database "postgresql://$(dbUser):$(dbPassword)@$(dbHost):$(dbPassword)/$(dbName)?sslmode=disable" -verbose down; \
	fi 

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test -v -cover ./...

# .PHONY: test-coverage
test-coverage:
#   go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out |
#   go tool cover -html coverage/coverage.out -o coverage/coverage.html |
#   open coverage/coverage.html 
# postgresql://sonar_user_checklos:S0N4RQUB3!@#$@154.12.237.18:31373/dh-user?sslmode=enable