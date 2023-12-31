.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

docker-rm-volume: ### Remove docker volume
	docker volume rm codespaces-blank_pg-data
.PHONY: docker-rm-volume

linter-golangci: ### Check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

swag: ### Generate swag docs
	swag init -g cmd/app/main.go
.PHONY: swag

test: ### Run test
	go test -v ./...
.PHONY: test

migrate-create:  ### Create new migration
	migrate create -ext sql -dir migrations 'EnrichInfo'
.PHONY: migrate-create