package logic

import uuid "github.com/satori/go.uuid"

// GetUUID 生成UUID
func GetUUID() string {
	return uuid.NewV4().String()
}
