.PHONY: migrate migrate_down migrate_up migrate_version docker prod docker_delve local swaggo test
VERSION ?= $(shell git describe --tags --always)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS ?= -X github.com/JamesHsu333/go-twitter/pkg/version.Version=$(VERSION) -X github.com/JamesHsu333/go-twitter/pkg/version.BuildDate=$(BUILD_DATE)
DOCKER_LOGIN ?= 
IMAGE_NAME ?= go-twitter
TAG ?=
MIGRATE ?=

# Main
run:
	go run ./cmd/api/main.go

build:
	go build -ldflags="$(LDFLAGS)" -o bin/ ./cmd/api/main.go

test:
	go test -cover ./...

# Build docker image
dockerfile:
	docker build -t $(DOCKER_LOGIN)/$(IMAGE_NAME):$(TAG) .

# Docker compose commands
local:
	echo "Starting local environment"
	docker-compose -f docker-compose.local.yml up -d --build

# Modules support
deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# Docker support
FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)

# Go migrate postgresql
force:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations force $(MIGRATE)

version:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations up $(MIGRATE)

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations down $(MIGRATE)

# Tools commands
linter:
	echo "Starting linters"
	golangci-lint run ./...

swaggo:
	echo "Starting swagger generating"
	swag init -g **/**/*.go
