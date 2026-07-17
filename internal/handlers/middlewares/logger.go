package middlewares

import (
	"fmt"
	"net/mail"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
)

type MaskFunc func(*[]byte)

var (
	// regist ฟังชัน replace ค่าใน request
	registeredMasks = []MaskFunc{}

	// ขนาด body response มากสุดที่จะแสดงใน log
	maxResponseLogSize = 4 * 1024 // 4 KB

	// ค่าที่ใช้ตรวจสอบคีย์ว่าเป็น password
	maskPasswordFields = [...]string{
		"password",
		"old_password",
		"new_password",
		"confirm_password",
		"pin",
		"passcode",
		"secret",
	}

	// ค่าที่ใช้ตรวจสอบคีย์ว่าเป็น token, key
	maskTokenFields = [...]string{
		"token",
		"access_token",
		"refresh_token",
		"apiKey",
		"session_id",
	}
)

// Logger logger
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		// เวลาก่อนเริ่ม process
		startTime := time.Now()

		// ส่งต่อให้ Handler หรือ Middleware ถัดไปทำงาน
		err := c.Next()
		if err != nil {
			return err
		}

		// คำนวณเวลาหลังจากประมวลผลเสร็จสิ้น
		latency := time.Since(startTime)

		userID := cctx.New(c).GetUserID()

		var b []byte
		if parameters := c.Locals(cctx.ParametersKey); parameters != nil {
			b, err = json.Marshal(parameters)
			if err != nil {
				b = []byte(`"<middlewares logger marshal error>"`)
			}

			// replace ค่า secert ใน request
			registerMask(maskPassword, maskToken, maskEmailSensitive)
			applyMasks(&b)
		}

		logs := logger.Log.With(
			zap.String("app", config.CF.App.ProjectID),
			zap.String("version", config.CF.App.Version),
			zap.String("host", c.Hostname()),
			zap.String("method", c.Method()),
			zap.String("route", c.Route().Path),
			zap.String("path", c.Path()),
			zap.String("request_id", c.RequestID()),
			zap.String("client_ip", getIP(c)),
			zap.String("user_agent", c.Get(fiber.HeaderUserAgent)),
			zap.Int("request_bytes", len(c.Request().Body())),
		)

		if len(b) > 0 {
			logs = logs.With(
				zap.Any("parameters", json.RawMessage(b)),
			)
		}
		if userID != "" {
			logs = logs.With(
				zap.String("user_id", userID),
			)
		}

		if c.OriginalURL() == fmt.Sprintf("%s/healthcheck", config.CF.App.ApiBaseUrl) {
			return nil
		}
		if strings.HasPrefix(c.OriginalURL(), fmt.Sprintf("%s/swagger", config.CF.Swagger.BaseURL)) {
			return nil
		}

		logs.Infow(
			"request completed",
			"latency_ms", latency.Milliseconds(), // หน่วย ms เป็นหน่วยที่นิยมที่สุดสำหรับ API Monitoring
			"response_status", c.Response().StatusCode(),
			"response_bytes", len(c.Response().Body()),
			"response", prepareResponseLog(c.Response().Body()),
		)

		return nil
	}
}

func registerMask(fn ...MaskFunc) {
	registeredMasks = append(registeredMasks, fn...)
}

func applyMasks(data *[]byte) {
	for _, fn := range registeredMasks {
		fn(data)
	}
}

func maskReplaceAll(b []byte, fields []string, mask string) []byte {
	for _, f := range fields {
		if !gjson.GetBytes(b, f).Exists() {
			continue
		}

		if bb, err := sjson.SetBytes(b, f, mask); err == nil {
			b = bb
		}
	}

	return b
}

func maskReplaceMiddle(b []byte, fields []string, head, tail int) []byte {
	for _, f := range fields {
		value := gjson.GetBytes(b, f)
		if !value.Exists() {
			continue
		}

		masked := maskMiddle(value.String(), head, tail)
		if bb, err := sjson.SetBytes(b, f, masked); err == nil {
			b = bb
		}
	}

	return b
}

func maskMiddle(s string, head, tail int) string {
	if s == "" {
		return ""
	}

	if head < 0 {
		head = 0
	}

	if tail < 0 {
		tail = 0
	}

	if len(s) <= head+tail {
		return "****"
	}

	return s[:head] + "..." + s[len(s)-tail:]
}

// maskPassword mask password
// แทนค่ารหัสป็น **********
// ตย. pass1234 >> **********
func maskPassword(b *[]byte) {
	*b = maskReplaceAll(*b, maskPasswordFields[:], "**********")
}

func maskToken(b *[]byte) {
	*b = maskReplaceMiddle(*b, maskTokenFields[:], 8, 4)
}

func maskEmailSensitive(b *[]byte) {
	var v any
	if err := json.Unmarshal(*b, &v); err != nil {
		return
	}
	maskEmail(v)

	bb, err := json.Marshal(v)
	if err != nil {
		return
	}

	*b = bb
}

func maskEmailString(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	name := parts[0]
	domain := parts[1]

	if len(name) <= 2 {
		return "**@" + domain
	}

	return name[:2] + strings.Repeat("*", len(name)-2) + "@" + domain
}

// maskEmail mask email
// แทนค่าอีเมล์เป็น **
// ตย. user@mail.com >> us**@mail.com
func maskEmail(v any) {
	switch x := v.(type) {
	case map[string]any:
		for k, val := range x {
			switch vv := val.(type) {
			case string:
				if _, err := mail.ParseAddress(vv); err == nil {
					x[k] = maskEmailString(vv)
				}
			default:
				maskEmail(vv)
			}
		}

	case []any:
		for i, val := range x {
			switch vv := val.(type) {
			case string:
				if _, err := mail.ParseAddress(vv); err == nil {
					x[i] = maskEmailString(vv)
				}
			default:
				maskEmail(vv)
			}
		}
	}
}

// getIP get client ip
func getIP(c fiber.Ctx) string {
	var ips []string
	if c.IP() != "" {
		ips = append(ips, c.IP())
	}
	if len(c.IPs()) > 0 {
		ips = append(ips, c.IPs()...)
	}

	return strings.Join(ips, ", ")
}

// prepareResponseLog prepare response log
// เช็ค และแสดงค่า body
// ถ้าขนาด body มากกว่า max จะแสดงเป็น truncated=true, size=bodySize
func prepareResponseLog(b []byte) any {
	if len(b) <= maxResponseLogSize {
		// replace ค่า secert
		maskPassword(&b)
		maskToken(&b)
		maskEmailSensitive(&b)

		var v any
		if json.Unmarshal(b, &v) == nil {
			return v
		}
	}

	return map[string]any{
		"truncated": true,
		"size":      len(b),
	}
}
