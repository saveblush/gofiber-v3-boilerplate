package utils

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
)

var (
	patternKey      = "%s-%s"
	patternKeyMode  = "%s-%s-%s"
	keySessionLogin = "sessionlogin"
)

// Pointer pointer
func Pointer[Value any](v Value) *Value {
	return &v
}

// UUID
func UUID() string {
	v7, err := uuid.NewRandom()
	if err != nil {
		return ""
	}

	return v7.String()
}

// SetKey set key
// set key สำหรับเก็บ cache
func SetKey(sign, user string) string {
	return fmt.Sprintf(patternKeyMode, config.CF.App.ProjectId, sign, user)
}

// SetKeySessionLogin set key session login
// set key สำหรับ cache session login
// ตย. system-session-userid-sessionid
func SetKeySessionLogin(userID, sessionID string) string {
	return SetKey(keySessionLogin, fmt.Sprintf(patternKey, userID, sessionID))
}
