#!/bin/bash

set -e

PROJECT_NAME="ServerManagerSystem"
BINARY="main"
Version=$(date "+%Y%m%d%H%M%S")
OUTPUT_DIR=output

APP_NAME=${PROJECT_NAME}
APP_VERSION=SM${Version}
BUILD_VERSION=$(git log -1 --oneline)
BUILD_TIME=$(date "+%FT%T%z")
GIT_REVISION=$(git rev-parse --short HEAD)
GIT_BRANCH=$(git name-rev --name-only HEAD)
GO_VERSION=$(go version)
SYS_VERSION=$(uname -m)

CGO_ENABLED=1 go build -a -installsuffix cgo -v -mod=vendor \
  -ldflags "-s -w -X 'github.com/uouuou/ServerManagerSystem/models.AppName=${APP_NAME}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.AppVersion=${APP_VERSION}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.BuildVersion=${BUILD_VERSION}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.BuildTime=${BUILD_TIME}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.GitRevision=${GIT_REVISION}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.GitBranch=${GIT_BRANCH}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.GoVersion=${GO_VERSION}'" \
  -o ${OUTPUT_DIR}/${PROJECT_NAME}_${SYS_VERSION} ${BINARY}.go

#删除原有的UPX压缩包并创建新的压缩包
rm -rf ${OUTPUT_DIR}/${PROJECT_NAME}_${SYS_VERSION}_upx
upx -5 -o ${OUTPUT_DIR}/${PROJECT_NAME}_${SYS_VERSION}_upx ${OUTPUT_DIR}/${PROJECT_NAME}_${SYS_VERSION}

# 执行结束后打印版本号
echo "--------------------------"
echo "--------------------------"
echo "         版本号            "
echo "${APP_VERSION}"
echo "--------------------------"
echo "--------------------------"
