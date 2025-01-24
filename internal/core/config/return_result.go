package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
)

// RR -> for use to return result model
var (
	RR = &ReturnResult{}
)

// Language language
type Language string

const (
	// LanguageTH th
	LanguageTH Language = "th"
	// LanguageEN en
	LanguageEN Language = "en"
)

// IsValid is valid
func (l Language) IsValid() bool {
	return l == LanguageTH || l == LanguageEN
}

// String string
func (l Language) String() string {
	return string(l)
}

// Result result
type Result struct {
	Code        int               `json:"code" mapstructure:"code"`
	Description LocaleDescription `json:"message" mapstructure:"localization"`
}

// SwaggerInfoResult swagger info result
type SwaggerInfoResult struct {
	Code        int    `json:"code"`
	Description string `json:"message"`
}

// WithLocale with locale
func (rs Result) WithLocale(c fiber.Ctx) Result {
	locale, ok := c.Locals("lang").(Language)
	if !ok {
		rs.Description.Locale = LanguageTH
	}
	rs.Description.Locale = locale

	return rs
}

// Error error description
func (rs Result) Error() string {
	if rs.Description.Locale == LanguageTH {
		return rs.Description.TH
	}

	return rs.Description.EN
}

// ErrorCode error code
func (rs Result) ErrorCode() int {
	return rs.Code
}

// HTTPStatusCode http status code
func (rs Result) HTTPStatusCode() int {
	switch rs.Code {
	case 200: // success
		return fiber.StatusOK
	case 400: // bad request
		return fiber.StatusBadRequest
	case 404: // connection_error
		return fiber.StatusNotFound
	case 401: // unauthorized
		return fiber.StatusUnauthorized
	case 403: // forbidden
		return fiber.StatusForbidden
	}

	return fiber.StatusInternalServerError
}

// ReturnResult return result model
type ReturnResult struct {
	JSONDuplicateOrInvalidFormat Result `mapstructure:"json_duplicate_or_invalid_format"`
	InvalidToken                 Result `mapstructure:"invalid_token"`
	InvalidPermissionRole        Result `mapstructure:"invalid_permission_role"`
	TokenNotFound                Result `mapstructure:"token_not_found"`
	UserNotFound                 Result `mapstructure:"user_not_found"`
	EmployeeNotFound             Result `mapstructure:"employee_not_found"`

	Internal struct {
		Success          Result `mapstructure:"success"`
		General          Result `mapstructure:"general"`
		BadRequest       Result `mapstructure:"bad_request"`
		ConnectionError  Result `mapstructure:"connection_error"`
		DatabaseNotFound Result `mapstructure:"database_not_found"`
		Unauthorized     Result `mapstructure:"unauthorized"`
		Forbidden        Result `mapstructure:"forbidden"`
		TooManyRequests  Result `mapstructure:"too_many_requests"`
	} `mapstructure:"internal"`
}

// LocaleDescription locale description
type LocaleDescription struct {
	EN     string   `mapstructure:"en"`
	TH     string   `mapstructure:"th"`
	Locale Language `mapstructure:"success"`
}

// MarshalJSON marshall json
func (ld LocaleDescription) MarshalJSON() ([]byte, error) {
	if ld.Locale == LanguageTH {
		return json.Marshal(&ld.TH)
	}

	return json.Marshal(&ld.EN)
}

// UnmarshalJSON unmarshal json
func (ld *LocaleDescription) UnmarshalJSON(data []byte) error {
	var res string
	err := json.Unmarshal(data, &res)
	if err != nil {
		return err
	}
	ld.EN = res
	ld.Locale = LanguageEN

	return nil
}

// InitReturnResult init return result
func InitReturnResult() error {
	v := viper.New()
	v.AddConfigPath(filePath)
	v.SetConfigName(fileNameReturnResult)
	v.SetConfigType(fileExtension)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logger.Log.Errorf("read config file error: %s", err)
		return err
	}

	if err := bindingReturnResult(v, RR); err != nil {
		logger.Log.Errorf("binding config error: %s", err)
		return err
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		logger.Log.Infof("config file changed: %s", e.Name)
		if err := v.Unmarshal(&CF); err != nil {
			logger.Log.Errorf("binding config error: %s", err)
		}
	})
	v.WatchConfig()

	return nil
}

// bindingReturnResult binding return result
func bindingReturnResult(vp *viper.Viper, rr *ReturnResult) error {
	if err := vp.Unmarshal(&rr); err != nil {
		logger.Log.Errorf("unmarshal config error: %s", err)
		return err
	}

	return nil
}

// CustomMessage custom message
func (rr *ReturnResult) CustomMessage(messageEN, messageTH string, code ...int) *Result {
	result := &Result{
		Code: 999,
		Description: LocaleDescription{
			EN: messageEN,
			TH: messageTH,
		},
	}
	if code != nil {
		result.Code = code[0]
	}

	return result
}
