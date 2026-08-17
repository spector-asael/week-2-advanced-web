# CMPS4191 - Advanced Web Development
# Asael Tobar
# August 16th, 2026

include .envrc 

.PHONY: run 
run: 
	@echo "Starting the application..."
	@go run ./cmd/api -port=${PORT} -db-dsn=${DB_DSN} -env=development -report-delay="5s"

## db/psql: Connect to the banking database using psql
.PHONY: db
db:
	psql ${DB_DSN}


## db/migrations/new name=$1: Create a new database migration
.PHONY: migrations/new
migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## migrations/up: Apply all up database migrations
.PHONY: migrations/up
migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_DSN} up

## db/migrations/down: Revert all migrations
.PHONY: migrations/down
migrations/down:
	@echo 'Reverting all migrations...'
	migrate -path ./migrations -database ${DB_DSN} down

## db/migrations/fix version=$1: Force schema_migrations version
.PHONY: migrations/fix
migrations/fix:
	@echo 'Forcing schema migrations version to ${version}...'
	migrate -path ./migrations -database ${DB_DSN} force ${version}

.PHONY: createjob
createjob:
	@echo "Creating a new job..."
	@curl -X POST -H "Content-Type: application/json" -d '{"job_type":"example","status":"pending","payload":"{\"data\":\"example\"}","result":""}' http://localhost:${PORT}/jobs

