include .env
export

all: build test

run: run-user run-proxy

run-user:
	cd cmd/user && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/Watch2Gather/server/cmd/user
.PHONY: run-product

run-proxy:
	cd cmd/proxy && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/Watch2Gather/server/cmd/proxy
.PHONY: run-proxy

docker-compose: docker-compose-stop docker-compose-start
.PHONY: docker-compose

docker-compose-start:
	docker-compose up --build
.PHONY: docker-compose-start

docker-compose-stop:
	docker-compose down --remove-orphans -v
.PHONY: docker-compose-stop

docker-compose-core: docker-compose-core-stop docker-compose-core-start

docker-compose-core-start:
	docker-compose -f docker-compose-core.yaml up --build
.PHONY: docker-compose-core-start

docker-compose-core-stop:
	docker-compose -f docker-compose-core.yaml down --remove-orphans -v
.PHONY: docker-compose-core-stop

docker-compose-build:
	docker-compose down --remove-orphans -v
	docker-compose build
.PHONY: docker-compose-build

wire:
	cd internal/user/app && wire && cd -
	# cd internal/app1/app && wire && cd - && \
	# cd internal/app2/app && wire && cd - && \
	# cd internal/app3/app && wire && cd -
.PHONY: wire

sqlc:
	sqlc generate
.PHONY: sqlc

buf:
	buf generate
.PHONY: buf

test:
	go test -v main.go

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

clean:
	go clean
