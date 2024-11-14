package cctx

import (
	"fmt"
	"reflect"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
)

// Localizer localizer
type Localizer interface {
	WithLocale(c *Context)
	GetLanguage() config.Language
}

// Localization localization
func (c *Context) Localization(i interface{}, depth int) {
	const (
		key = "Localization"
	)
	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice:
		s := reflect.ValueOf(i)
		for i := 0; i < s.Len(); i++ {
			c.Localization(s.Index(i).Interface(), depth)
		}

	default:
		formValue := reflect.ValueOf(i)
		if !generic.Equal(formValue.Kind(), reflect.Map) {
			localizer, ok := formValue.Interface().(Localizer)
			if ok {
				localizer.WithLocale(c)
			}
			if generic.Equal(formValue.Kind(), reflect.Ptr) {
				formValue = formValue.Elem()
			}
			if generic.Equal(formValue.Kind(), reflect.Struct) {
				t := reflect.TypeOf(formValue.Interface())
				for i := 0; i < t.NumField(); i++ {
					fieldName := t.Field(i).Name
					value := formValue.FieldByName(fieldName)
					if value.IsValid() {
						if generic.Equal(fieldName, key) && generic.Equal(value.Kind(), reflect.Struct) {
							t := reflect.TypeOf(value.Interface())
							for i := 0; i < t.NumField(); i++ {
								fieldName := t.Field(i).Name
								value := formValue.FieldByName(fieldName)
								if value.IsValid() {
									switch localizer.GetLanguage() {
									case "th":
										v := formValue.FieldByName(fmt.Sprintf("%sTH", fieldName))
										if v.IsValid() {
											value.Set(v)
										}
									default:
										v := formValue.FieldByName(fmt.Sprintf("%sEN", fieldName))
										if v.IsValid() {
											value.Set(v)
										}
									}
								}
							}
						}

						_, ok := value.Interface().(Localizer)
						if ok {
							if value.IsValid() && !value.IsNil() {
								c.Localization(value.Interface(), depth+1)
							}
						}

						if depth <= compositeFormDepth && (generic.Equal(value.Kind(), reflect.Slice) || fieldName == entities) {
							c.Localization(value.Interface(), depth+1)
						}
					}
				}
			}
		}
	}
}
