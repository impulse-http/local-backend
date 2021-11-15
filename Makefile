SHELL := /usr/bin/env bash

.PHONY: database
database:
	@echo "Initializing database"
	cd migrations; goose sqlite3 ../dist/impulse.db up

.PHONY: start-dev
start-dev:
	@echo "Starting development server..."
	make database
	go build -o dist/service cmd/main.go
	./dist/service

.PHONY: deps
deps:
	@echo "Getting dependencies..."
	go get -u github.com/pressly/goose/v3/cmd/goose
