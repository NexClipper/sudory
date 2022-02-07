#!/bin/sh

initdb() {
	/app/init-db.sh ${SUDORY_DB_HOST} ${SUDORY_DB_PORT} ${SUDORY_DB_SCHEME} /app/sudory.sql ${SUDORY_DB_EXPORT_PATH} ${SUDORY_DB_ROOT_USERNAME} ${SUDORY_DB_ROOT_PASSWORD} ${SUDORY_DB_SERVER_USERNAME} ${SUDORY_DB_SERVER_PASSWORD}
}

apprun() {
	/app/sudory-server -config /app/sudory-server.yml
}

initdb
apprun