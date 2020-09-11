GOPATH = $(shell go env GOPATH)
GOBIN = $(GOPATH)/bin

############################################################################
### Custom Installs ########################################################
############################################################################
GO_MIGRATE = $(GOBIN)/migrate
$(GO_MIGRATE):
	@echo ">> Couldn't find go-migrate; installing..."
	go get -tags 'sqlite3' -u github.com/golang-migrate/migrate/v4/cmd/migrate


############################################################################
### Targets ################################################################
############################################################################

all: clean bot

MIGRATIONS = ./migrations
SQLITE_MIGRATIONS = $(MIGRATIONS)/sqlite
DATA = ./data
SQLITE_DB = $(DATA)/sqlite/pikamon.db

migrate-up:
	@echo "Migrating up..."
	@mkdir -p ./data/sqlite
	@migrate -path $(SQLITE_MIGRATIONS) -database sqlite3://$(SQLITE_DB) up $(SCHEMA_VERSION)

migrate-down:
	@echo "Migrating down..."
	@migrate -path $(SQLITE_MIGRATIONS) -database sqlite3://$(SQLITE_DB) down $(SCHEMA_VERSION)

bot:
	@echo "Building bot binary for use on local system..."
	@env CGO_ENABLED=1 go build -o bin/pikamon ./cmd/pikamon

clean:
	@echo "Cleaning bin/..."
	@rm -rf bin/*

project-utils: $(GO_MIGRATE)
	@echo "Installing project utilities..."

docker-image:
	@echo "Building docker image..."
	@docker build -t pikamon .
