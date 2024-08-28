ITOJU_BINARY=itojuApp
# ==================================================================================== # 
# HELPERS 
# ==================================================================================== #

## help: print this help message
.PHONY: help 
help: 
	@echo 'Usage:' 
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm 
confirm: 
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== # 
# DEVELOPMENT 
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo 'Starting server...'
	go run ./cmd/main

## db/psql: connect to the database using psql and docker
.PHONY: db/psql
db/psql:
	@echo 'connecting to db'
	docker exec -it post-db bash
	psql postgres://itojudb:itojudb@localhost/itojudb?sslmode=disable
 
## db/migrate/up: apply all up database migrations
.PHONY: db/migrate/up
db/migrate/up:
	echo 'Running up migrations...'
	@cd internal/sql/migrations/ && goose postgres postgres://itojudb:itojudb@localhost/itojudb up && goose postgres postgres://koyeb-adm:rcHo1Ck7BYmf@ep-tiny-mode-a2d0vyca.eu-central-1.pg.koyeb.app/Itoju-ky up && goose postgres postgres://djjsagev:WG11sRXwe2q1C0I9-3XhTZywTnhbZQPJ@stampy.db.elephantsql.com/djjsagev up

.PHONY: db/migrate/upt
db/migrate/upt:
	echo 'Running up migrations...'
	@cd internal/sql/migrations/ && goose postgres postgres://itojudb:itojudb@localhost/itojudb up 
.PHONY: db/migrate/downt
db/migrate/downt:
	@echo 'Running down migrations...'
	@cd internal/sql/migrations/ && goose postgres postgres://itojudb:itojudb@localhost/itojudb down

## db/migrate/down: apply all down database migrations
.PHONY: db/migrate/down
db/migrate/down:
	@echo 'Running down migrations...'
	@cd internal/sql/migrations/ && goose postgres postgres://itojudb:itojudb@localhost/itojudb down && goose postgres postgres://koyeb-adm:rcHo1Ck7BYmf@ep-tiny-mode-a2d0vyca.eu-central-1.pg.koyeb.app/Itoju-ky down && goose postgres postgres://djjsagev:WG11sRXwe2q1C0I9-3XhTZywTnhbZQPJ@stampy.db.elephantsql.com/djjsagev down

# ==================================================================================== # 
# QUALITY CONTROL 
# ==================================================================================== # 
## audit: tidy dependencies and format, vet and test all code 

.PHONY: audit 
audit: vendor
	@echo 'Formatting code...' 
	go fmt ./... 
	@echo 'Vetting code...' 
	go vet ./... 
	staticcheck ./... 
	@echo 'Running tests...' 
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies 
.PHONY: vendor 
vendor: 
	@echo 'Tidying and verifying module dependencies...' 
	go mod tidy 
	go mod verify 
	@echo 'Vendoring dependencies...' 
	go mod vendor

# ==================================================================================== # 
# BUILD 
# ==================================================================================== # 

## build/api: build the cmd/api application 
.PHONY: build/api 
build/api: 
	@echo 'Building cmd/api...' 
	env GOOS=linux CGO_ENABLED=0 go build -o bin/${ITOJU_BINARY} ./cmd/main
	# go build -ldflags='-s' -o=./bin/api ${ITOJU_BINARY} ./cmd/main
	# GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## build/docker: build the docker application 
.PHONY: build/docker 
build/docker: build/api
	@echo 'Building docker...' 
	docker build -t itojuapp . 

.PHONY: run/docker 
run/docker: build/docker
	@echo 'Building docker...' 
	docker run -e DB_URL=postgres://djjsagev:WG11sRXwe2q1C0I9-3XhTZywTnhbZQPJ@stampy.db.elephantsql.com/djjsagev itojuapp
