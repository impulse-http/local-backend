SHELL := /usr/bin/env bash

.PHONY: initdb
initdb:
	@echo "Initializing database..."
	mkdir -p dist
	cd migrations; goose sqlite3 ../dist/impulse.db up

.PHONY: dropdb
dropdb:
	@echo "Dropping database..."
	rm -f dist/impulse.db

.PHONY: reinitdb
reinitdb:
	@echo "Reinitializing database..."
	make dropdb
	make initdb

.PHONY: start-dev
start-dev:
	@echo "Starting development server..."
	make initdb
	go build -o dist/service cmd/main.go
	./dist/service

.PHONY: deps
deps:
	@echo "Getting dependencies..."
	go get -u github.com/pressly/goose/v3/cmd/goose
