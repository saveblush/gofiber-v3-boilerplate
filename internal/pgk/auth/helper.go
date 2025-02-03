package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// comparePassword compare password
func comparePassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// createCacheSessionLogin create cache session login
// เก็บ cache session login ของ token.AccessToken เพื่อใช้เช็คว่า token ยังมี session อยู่ในระบบ
// ใช้กับเส้น me และ refreshtoken
func (s *service) createCacheSessionLogin(userID, sessionID string, data *models.Token) error {
	key := utils.SetKeySessionLogin(userID, sessionID)
	fmt.Println(key)
	// ใช้ timeout ของ RefreshToken
	_ = s.cache.Set(key, data, s.config.JWT.RefreshExpireTime)

	return nil
}

// getCacheSessionLogin get cache session login
// ดึง cache session login
func (s *service) getCacheSessionLogin(userID, sessionID string) (*models.Token, error) {
	key := utils.SetKeySessionLogin(userID, sessionID)
	value := &models.Token{}
	_ = s.cache.Get(key, value)

	return value, nil
}

// clearCacheSessionLogin clear cache session login
// เคลีย cache session login
func (s *service) clearCacheSessionLogin(userID, sessionID string) error {
	key := utils.SetKeySessionLogin(userID, sessionID)
	_ = s.cache.Delete(key)

	return nil
}

// checkSessionLogin check session login
// เช็ค session login
func (s *service) checkSessionLogin(c *cctx.Context, userID, sessionID, tokenCheck string) error {
	if generic.IsEmpty(tokenCheck) {
		logger.Log.Errorf("token not found")
		return c.ErrorUnauthorized(s.result.TokenNotFound.WithLocale(c.Ctx).Error())
	}

	session, err := s.getCacheSessionLogin(userID, sessionID)
	if err != nil {
		logger.Log.Errorf("find session login error: %s", err)
		return s.result.Internal.Unauthorized
	}
	if generic.IsEmpty(session) {
		logger.Log.Errorf("session login not found")
		return s.result.Internal.Unauthorized
	}

	if session.AccessToken == tokenCheck || session.RefreshToken == tokenCheck {
		return nil
	}

	return s.result.Internal.Unauthorized
}

// getHeaderToken get header token
func (s *service) getHeaderToken(c *cctx.Context) string {
	header := strings.Split(c.Get(fiber.HeaderAuthorization), " ")
	auth := strings.TrimSpace(header[1])

	if generic.IsEmpty(auth) {
		return ""
	}

	return auth
}
