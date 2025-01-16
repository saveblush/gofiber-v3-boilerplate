package cctx

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
)

const (
	pathKey            = "path"
	queryKey           = "query"
	compositeFormDepth = 3
	entities           = "Entities"
	LangKey            = "lang"
	UserKey            = "user"
	ParametersKey      = "parameters"
)

// Context context
type Context struct {
	fiber.Ctx
}

// New new custom fiber context
func New(c fiber.Ctx) *Context {
	return &Context{c}
}

// BindValue bind value
func (c *Context) BindValue(i interface{}, validate bool) error {
	switch c.Method() {
	case fiber.MethodGet:
		_ = c.Bind().Query(i)
		//c.additionalQueryParser(i)

	default:
		_ = c.Bind().Body(i)
	}

	c.PathParser(i, 1)
	c.TrimSpace(i, 1)
	c.Locals(ParametersKey, i)

	if validate {
		err := c.Validate(i)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate validate
func (c *Context) Validate(i interface{}) error {
	err := config.CF.Validator.Struct(i)
	if err != nil {
		return config.RR.CustomMessage(err.Error(), err.Error()).WithLocale(c.Ctx)
	}

	return nil
}

// TrimSpace trim space
func (c *Context) TrimSpace(i interface{}, depth int) {
	e := reflect.ValueOf(i).Elem()
	for i := 0; i < e.NumField(); i++ {
		if depth <= compositeFormDepth && e.Type().Field(i).Type.Kind() == reflect.Struct {
			depth++
			c.TrimSpace(e.Field(i).Addr().Interface(), depth)
		}

		if e.Type().Field(i).Type.Kind() != reflect.String {
			continue
		}

		value := e.Field(i).String()
		e.Field(i).SetString(strings.TrimSpace(value))
	}
}

// PathParser parse path param
func (c *Context) PathParser(i interface{}, depth int) {
	formValue := reflect.ValueOf(i)
	if formValue.Kind() == reflect.Ptr {
		formValue = formValue.Elem()
	}

	t := reflect.TypeOf(formValue.Interface())
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		paramValue := formValue.FieldByName(fieldName)
		if paramValue.IsValid() {
			if depth < compositeFormDepth && paramValue.Kind() == reflect.Struct {
				depth++
				c.PathParser(paramValue.Addr().Interface(), depth)
			}
			tag := t.Field(i).Tag.Get(pathKey)
			if tag != "" {
				setValue(paramValue, c.Params(tag))
			}
		}
	}
}

func setValue(paramValue reflect.Value, value string) {
	if paramValue.IsValid() && value != "" {
		switch paramValue.Kind() {
		case reflect.Uint:
			number, _ := strconv.ParseUint(value, 10, 32)
			paramValue.SetUint(number)

		case reflect.String:
			paramValue.SetString(value)

		default:
			number, err := strconv.Atoi(value)
			if err != nil {
				paramValue.SetString(value)
			} else {
				paramValue.SetInt(int64(number))
			}
		}
	}
}

/*func (c *Context) additionalQueryParser(i interface{}) {
	formValue := reflect.ValueOf(i)
	if formValue.Kind() == reflect.Ptr {
		formValue = formValue.Elem()
	}

	t := reflect.TypeOf(formValue.Interface())
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		value := formValue.FieldByName(fieldName)
		if generic.Equal(value.Kind(), reflect.Struct) {
			c.additionalQueryParser(value.Addr().Interface())
		}
		if generic.Equal(value.Kind(), reflect.Slice) {
			if value.IsValid() {
				for i := 0; i < value.Len(); i++ {
					if value.Len() == 1 && generic.Equal(value.Index(i).Kind(), reflect.String) {
						v := value.Index(i).Interface()
						str := strings.Split(fmt.Sprintf("%v", v), ",")
						if t.Field(i).Tag.Get(queryKey) != "" {
							value.Set(reflect.ValueOf(str))
						}
					} else if generic.Equal(value.Index(i).Kind(), reflect.Struct) {
						c.additionalQueryParser(value.Index(i).Addr().Interface())
					}
				}
			}
		}
	}
}*/

// GetLanguage get language
func (c *Context) GetLanguage() config.Language {
	locale, ok := c.Locals("lang").(string)
	if !ok {
		return config.LanguageTH
	}
	if locale != "en" {
		return config.LanguageTH
	}

	return config.Language(locale)
}

// GetClientIP get client ip
func (c *Context) GetClientIP() string {
	var ip string
	for _, ips := range c.IPs() {
		ip = ips
	}
	if ip == "" {
		ip = c.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.Get("X-Client-IP")
	}
	if ip == "" {
		ip = c.Get("X-Real-Ip")
	}

	return ip
}

// GetClientUserAgent get client user agent
func (c *Context) GetClientUserAgent() string {
	return c.Get(fiber.HeaderUserAgent)
}
