package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"strings"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
	"github.com/golang-jwt/jwt/v5"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

var (
	prefixBasicKey = "basic_key_"
	prefixAdminKey = "admin_key_"
	prefixApiKey   = "api_key_"
	prefixWsKey    = "ws_key_"
)

// AuthorizationRequired authorization jwt and basicauth
func AuthorizationRequired() fiber.Handler {
	users := make(map[string]string)
	for _, item := range config.CF.App.Sources {
		if strings.HasPrefix(item.User, prefixBasicKey) {
			// ตัด text basic_key_
			users[item.User[len(prefixBasicKey):]] = item.Password
		}
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
		Extractor:  extractors.FromAuthHeader("Bearer"),
		KeyFunc: func(t *jwt.Token) (any, error) {
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
			if strings.HasPrefix(item.User, prefixAdminKey) {
				// ตัด text admin_key_
				users[item.User[len(prefixAdminKey):]] = item.Password
			}
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
			Extractor: extractors.FromHeader("X-API-Key"),
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

// validateAPIKey verify auth key
func validateAPIKey(_ fiber.Ctx, prefixKey, key string) (bool, error) {
	if prefixKey == "" {
		return false, errors.New("missing Prefix Key")
	}

	keys := make(map[string]string)
	for _, item := range config.CF.App.Sources {
		if strings.HasPrefix(item.User, prefixKey) {
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

// ValidateAPIKey verify api-key
func ValidateAPIKey(c fiber.Ctx, key string) (bool, error) {
	return validateAPIKey(c, prefixApiKey, key)
}

// ValidateWebsocketKey verify websocket-key
func ValidateWebsocketKey(c fiber.Ctx, key string) (bool, error) {
	return validateAPIKey(c, prefixWsKey, key)
}
