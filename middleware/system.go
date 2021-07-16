package middleware

import (
	"os"
	"runtime"
)

// GetArchitecture 获取操作系统制式arm64还是amd64
func GetArchitecture() (architecture string, system string) {
	architecture = runtime.GOARCH
	system = runtime.GOOS
	return
}

// FileExist 获取某个文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
