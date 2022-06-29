#!/usr/bin/env bash

MIGRATE_DRIVER="mysql"

go get -d -u github.com/golang-migrate/migrate/v4

go install -tags $MIGRATE_DRIVER github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate -version

read -p "Press Enter to Continue"