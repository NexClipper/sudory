#!/usr/bin/env bash

pwd=$PWD



# Show env vars
grep -v '^#' .sudory.env

# Export env vars
export $(grep -v '^#' .sudory.env | xargs)

path="$(dirname "$0")"
echo $pwd
echo $path

cd ../sudory

../cmd/repl.sh


cd $pwd