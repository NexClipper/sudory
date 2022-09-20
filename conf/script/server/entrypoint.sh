#!/bin/bash

apprun() {
	exec /app/sudory-server -config '/app/conf/sudory-server.yml'
}

apprun
