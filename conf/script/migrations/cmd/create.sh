#!/usr/bin/env bash

read -p "Enter migrate create name : " NAME

migrate create -seq -dir ../ -ext sql "$NAME"

# read -p "Press Enter to Continue"