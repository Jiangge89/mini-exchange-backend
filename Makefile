.PHONY: up down build test lint clean

up:
	docker compose up -d

down:
	docker compose down

build: build-gateway build-indexer build-txservice build-contracts

build-gateway:
	cd gateway && go build ./...

build-indexer:
	cd indexer && go build ./...

build-txservice:
	cd txservice && ./gradlew build -x test

build-contracts:
	cd contracts && forge build

test: test-gateway test-indexer test-txservice test-contracts

test-gateway:
	cd gateway && go test ./...

test-indexer:
	cd indexer && go test ./...

test-txservice:
	cd txservice && ./gradlew test

test-contracts:
	cd contracts && forge test

lint: lint-gateway lint-indexer

lint-gateway:
	cd gateway && go vet ./...

lint-indexer:
	cd indexer && go vet ./...

clean:
	cd gateway && go clean ./...
	cd indexer && go clean ./...
	cd txservice && ./gradlew clean
	cd contracts && forge clean
