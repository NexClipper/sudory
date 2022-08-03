#!/bin/sh

if [[ -f /vault/secrets/ncdb-admin-login ]]; then
    source /vault/secrets/ncdb-admin-login
    export SUDORY_DB_ROOT_USERNAME="${SUDORY_DB_ROOT_USERNAME}"
    export SUDORY_DB_ROOT_PASSWORD="${SUDORY_DB_ROOT_PASSWORD}"
    export SUDORY_DB_SERVER_USERNAME="${SUDORY_DB_SERVER_USERNAME}"
    export SUDORY_DB_SERVER_PASSWORD="${SUDORY_DB_SERVER_PASSWORD}"
fi

initdb() {
	./init-db.sh $SUDORY_DB_HOST $SUDORY_DB_PORT $SUDORY_DB_SCHEME /app/sudory.sql $SUDORY_DB_EXPORT_PATH $SUDORY_DB_ROOT_USERNAME $SUDORY_DB_ROOT_PASSWORD $SUDORY_DB_SERVER_USERNAME $SUDORY_DB_SERVER_PASSWORD
}

apprun() {
	/app/sudory-server -config '/app/conf/sudory-server.yml'
}

initdb
apprun

