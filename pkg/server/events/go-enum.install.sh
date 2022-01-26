#!/bin/bash

##
## go-enum 설치를 위한 스크립트
##
## 불필요한 의존성 정보를 남기지 않기 위해
## export GO111MODULE=off #설정
##
export GO111MODULE=off
go get -v github.com/abice/go-enum
export GO111MODULE=on



##
## 코드 자동 생성 명령
## $ go-enum --file=<go-enum_코드파일이름.go> --names --nocase=true
##