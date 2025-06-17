package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
	"github.com/golang-jwt/jwt/v5"
	jwtware "github.com/saveblush/gofiber3-contrib/jwt"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// AuthorizationRequired authorization jwt and basicauth
func AuthorizationRequired() fiber.Handler {
	users := make(map[string]string)
	for _, item := range config.CF.App.Sources {
		users[item.User] = item.Password
	}

	basicAuth := basicauth.New(basicauth.Config{
		Users: users,
		Unauthorized: func(c fiber.Ctx) error {
			logger.Log.Error("authorization error: unauthorized")
			return fiber.NewError(config.RR.Internal.Unauthorized.HTTPStatusCode(), config.RR.InvalidToken.WithLocale(c).Error())
		},
	})

	return jwtware.New(jwtware.Config{
		Claims:     &models.TokenClaims{},
		SigningKey: jwtware.SigningKey{Key: []byte(config.CF.JWT.AccessSecretKey)},
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(config.CF.JWT.AccessSecretKey), nil
		},
		SuccessHandler: func(c fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return basicAuth(c)
		},
	})
}

// AuthorizationAdminRequired authorization admin basicauth
func AuthorizationAdminRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		users := make(map[string]string)
		for _, item := range config.CF.App.Sources {
			users[item.User] = item.Password
		}

		basicAuth := basicauth.New(basicauth.Config{
			Users: users,
			Unauthorized: func(c fiber.Ctx) error {
				logger.Log.Error("authorization admin error: unauthorized")
				return fiber.ErrUnauthorized
			},
		})

		return basicAuth(c)
	}
}

// AuthorizationAPIKey authorization x-api-key
func AuthorizationAPIKey() fiber.Handler {
	return func(c fiber.Ctx) error {
		auth := keyauth.New(keyauth.Config{
			KeyLookup: "header:x-api-key",
			Validator: func(c fiber.Ctx, key string) (bool, error) {
				return ValidateAPIKey(c, key)
			},
			SuccessHandler: func(c fiber.Ctx) error {
				return c.Next()
			},
			ErrorHandler: func(c fiber.Ctx, err error) error {
				logger.Log.Error("authorization x-api-key error: unauthorized")
				if err == keyauth.ErrMissingOrMalformedAPIKey {
					return fiber.NewError(config.RR.Internal.Unauthorized.HTTPStatusCode(), config.RR.InvalidToken.WithLocale(c).Error())
				}
				return fiber.ErrUnauthorized
			},
		})

		return auth(c)
	}
}

// ValidateAPIKey verify api-key
func ValidateAPIKey(c fiber.Ctx, key string) (bool, error) {
	keys := make(map[string]string)
	for _, item := range config.CF.App.Sources {
		if strings.HasPrefix(item.User, "api_key_") {
			keys[item.Password] = item.Password
		}
	}

	sourceKey, ok := keys[key]
	if !ok {
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}

	hashSourceKey := sha256.Sum256([]byte(sourceKey))
	hashKey := sha256.Sum256([]byte(key))
	if subtle.ConstantTimeCompare(hashSourceKey[:], hashKey[:]) == 1 {
		return true, nil
	}

	return false, keyauth.ErrMissingOrMalformedAPIKey
}
