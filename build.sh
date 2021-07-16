#!/bin/bash

set -e

PROJECT_NAME="ServerManagerSystem"
BINARY="main"
Version=$(date "+%Y%m%d%H%M%S")
OUTPUT_DIR=output
GOOS=$(go env GOOS)

APP_NAME=${PROJECT_NAME}
APP_VERSION=SM${Version}
BUILD_VERSION=$(git log -1 --oneline)
BUILD_TIME=$(date "+%FT%T%z")
GIT_REVISION=$(git rev-parse --short HEAD)
GIT_BRANCH=$(git name-rev --name-only HEAD)
GO_VERSION=$(go version)

# 构建一个ARM V8版本由于使用了SQLLite需要引入一个gcc编译器
GOARCH=arm64 GOOS=linux CGO_ENABLED=1 CC=/usr/lib/gcc-arm-10.2-2020.11-x86_64-aarch64-none-linux-gnu/bin/aarch64-none-linux-gnu-gcc-10.2.1 go build -a -installsuffix cgo -v -mod=vendor \
  -ldflags "-s -w -X 'github.com/uouuou/ServerManagerSystem/models.AppName=${APP_NAME}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.AppVersion=${APP_VERSION}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.BuildVersion=${BUILD_VERSION}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.BuildTime=${BUILD_TIME}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.GitRevision=${GIT_REVISION}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.GitBranch=${GIT_BRANCH}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.GoVersion=${GO_VERSION}'" \
  -o ${OUTPUT_DIR}/${PROJECT_NAME}_linux_armv8 ${BINARY}.go

# 构建一个x86_64版本由于sqllite的原因需要使用CGO_ENABLED=1
#CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=/usr/bin/x86_64-w64-mingw32-gcc-win32 go build -a -installsuffix cgo -v -mod=vendor \
#  -ldflags "-s -w -X 'github.com/uouuou/ServerManagerSystem/models.AppName=${APP_NAME}' \
#            -X 'github.com/uouuou/ServerManagerSystem/models.AppVersion=${APP_VERSION}' \
#            -X 'github.com/uouuou/ServerManagerSystem/models.BuildVersion=${BUILD_VERSION}' \
#            -X 'github.com/uouuou/ServerManagerSystem/models.BuildTime=${BUILD_TIME}' \
#            -X 'github.com/uouuou/ServerManagerSystem/models.GitRevision=${GIT_REVISION}' \
#            -X 'github.com/uouuou/ServerManagerSystem/models.GitBranch=${GIT_BRANCH}' \
#            -X 'github.com/uouuou/ServerManagerSystem/models.GoVersion=${GO_VERSION}'" \
#  -o ${OUTPUT_DIR}/${PROJECT_NAME}_windows_x86_64 ${BINARY}.go

# 构建一个amd64版本由于sqllite的原因需要使用CGO_ENABLED=1
CGO_ENABLED=1 go build -a -installsuffix cgo -v -mod=vendor \
  -ldflags "-s -w -X 'github.com/uouuou/ServerManagerSystem/models.AppName=${APP_NAME}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.AppVersion=${APP_VERSION}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.BuildVersion=${BUILD_VERSION}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.BuildTime=${BUILD_TIME}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.GitRevision=${GIT_REVISION}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.GitBranch=${GIT_BRANCH}' \
            -X 'github.com/uouuou/ServerManagerSystem/models.GoVersion=${GO_VERSION}'" \
  -o ${OUTPUT_DIR}/${PROJECT_NAME}_linux_amd64 ${BINARY}.go

#删除原有的UPX压缩包并创建新的压缩包
rm -rf ${OUTPUT_DIR}/${PROJECT_NAME}_linux_armv8_upx
rm -rf ${OUTPUT_DIR}/${PROJECT_NAME}_linux_amd64_upx
upx -5 -o ${OUTPUT_DIR}/${PROJECT_NAME}_linux_armv8_upx ${OUTPUT_DIR}/${PROJECT_NAME}_linux_armv8
upx -5 -o ${OUTPUT_DIR}/${PROJECT_NAME}_linux_amd64_upx ${OUTPUT_DIR}/${PROJECT_NAME}_linux_amd64

# 执行结束后打印版本号
echo "--------------------------"
echo "--------------------------"
echo "         版本号            "
echo "${APP_VERSION}"
echo "--------------------------"
echo "--------------------------"
