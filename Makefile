APP_NAME?=geogame
IMAGE=$(APP_NAME)
VERSION?=snap-shot
GOCMD=CGO_ENABLED=0 GOOS=linux go
GOCMDTEST=CGO_ENABLED=0 go
ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: vendor
vendor:
	$(GOCMD) mod vendor

.PHONY: build
build:
	$(GOCMD) build -mod vendor -ldflags "-X main.serviceVersion=$(VERSION)" -o go-app $(ROOT_DIR)

.PHONY: image
image:
	@echo "Building $(IMAGE):$(VERSION)"
	@docker build -f Dockerfile --build-arg=VERSION=$(VERSION) -t $(IMAGE):$(VERSION) .

.PHONY: test_unit
test_unit:
	$(GOCMDTEST) test ./... -mod=vendor -count=1

.PHONY: up
up:
	@docker-compose -f docker-compose.yaml up
# 	@make migrate

.PHONY: down
down:
	@docker-compose -f docker-compose.yaml down  --remove-orphans --volumes

.PHONY: migrate
migrate:
	docker run -v $(ROOT_DIR)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgres://geogameuser:password@localhost:5432/geo_game_db?sslmode=disable" up
