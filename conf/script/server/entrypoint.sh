#!/bin/sh

initdb() {
	/app/init-db.sh ${ROOT_DB_USER} ${ROOT_DB_PASSWORD} ${DB_HOST} ${DB_PORT} ${SUDORY_DB_SCHEME} /app/sudory.sql ${SUDORY_DB_EXPORT_PATH} ${SUDORY_DB_USER} ${SUDORY_DB_PASSWORD}
}

apprun() {
	/app/sudory-server -config /app/sudory-server.yml
}

initdb
apprun