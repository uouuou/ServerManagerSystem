package middleware

import (
	"github.com/satori/go.uuid"
)

// GetUUID 获取uuid
func GetUUID() string {
	return uuid.NewV4().String()
}
