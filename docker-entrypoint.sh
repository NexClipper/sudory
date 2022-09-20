#!/bin/sh

if [[ -f /vault/secrets/ncdb-admin-login ]]; then
    source /vault/secrets/ncdb-admin-login
    export SUDORY_DB_SERVER_USERNAME="${SUDORY_DB_SERVER_USERNAME}"
    export SUDORY_DB_SERVER_PASSWORD="${SUDORY_DB_SERVER_PASSWORD}"
fi

apprun() {
	/app/sudory-server -config '/app/conf/sudory-server.yml'
}

apprun

