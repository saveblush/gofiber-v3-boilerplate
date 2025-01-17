package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bytedance/sonic"
	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/spf13/viper"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
)

var (
	CF = &Configs{}
)

var (
	filePath                           = "./configs"
	fileExtension                      = "yml"
	fileNameConfig                     = "config"
	fileNameReturnResult               = "return_result"
	fileNameConfigAvailableStatus      = "config_available_status.yml"
	fileNameConfigAvailableDescription = "config_available_description.yml"
	AvailableStatusOnline              = "online"
	AvailableStatusOffline             = "offline"
)

// Environment environment
type Environment string

const (
	Develop    Environment = "develop"
	Production Environment = "prod"
)

// Production check is production
func (e Environment) Production() bool {
	return e == Production
}

type AvailableConfig struct {
	Status string `json:"status"`
}

type DatabaseConfig struct {
	Host         string        `mapstructure:"HOST"`
	Port         int           `mapstructure:"PORT"`
	Username     string        `mapstructure:"USERNAME"`
	Password     string        `mapstructure:"PASSWORD"`
	DatabaseName string        `mapstructure:"DATABASE_NAME"`
	DriverName   string        `mapstructure:"DRIVER_NAME"`
	Charset      string        `mapstructure:"CHARSET"`
	Timeout      string        `mapstructure:"TIMEOUT"`
	MaxIdleConns int           `mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns int           `mapstructure:"MAX_OPEN_CONNS"`
	MaxLifetime  time.Duration `mapstructure:"MAX_LIFE_TIME"`
}

type UserPassConfig struct {
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
}

type Configs struct {
	UniversalTranslator *ut.UniversalTranslator
	Validator           *validator.Validate
	Path                string

	App struct {
		AvailableStatus string           // สถานะปิด/เปิดระบบ [on/off]
		ProjectId       string           `mapstructure:"PROJECT_ID"`
		ProjectName     string           `mapstructure:"PROJECT_NAME"`
		Version         string           `mapstructure:"VERSION"`
		WebBaseUrl      string           `mapstructure:"WEB_BASE_URL"`
		ApiBaseUrl      string           `mapstructure:"API_BASE_URL"`
		Port            int              `mapstructure:"PORT"`
		Environment     Environment      `mapstructure:"ENVIRONMENT"`
		Issuer          string           `mapstructure:"ISSUER"`
		Sources         []UserPassConfig `mapstructure:"SOURCES"`
	} `mapstructure:"APP"`

	HTTPServer struct {
		Prefork   bool `mapstructure:"PREFORK"`
		RateLimit struct {
			Max        int           `mapstructure:"MAX"`
			Expiration time.Duration `mapstructure:"EXPIRATION"`
			Enable     bool          `mapstructure:"ENABLE"`
		} `mapstructure:"RATELIMIT"`
	} `mapstructure:"HTTP_SERVER"`

	Web struct {
		DateFormat     string `mapstructure:"DATE_FORMAT"`
		DateTimeFormat string `mapstructure:"DATETIME_FORMAT"`
		TimeFormat     string `mapstructure:"TIME_FORMAT"`
	} `mapstructure:"WEB"`

	Swagger struct {
		Title       string `mapstructure:"TITLE"`
		Version     string `mapstructure:"VERSION"`
		Host        string `mapstructure:"HOST"`
		BaseURL     string `mapstructure:"BASE_URL"`
		Description string `mapstructure:"DESCRIPTION"`
		Enable      bool   `mapstructure:"ENABLE"`
	} `mapstructure:"SWAGGER"`

	JWT struct {
		AccessSecretKey   string        `mapstructure:"ACCESS_SECRET_KEY"`
		RefreshSecretKey  string        `mapstructure:"REFRESH_SECRET_KEY"`
		AccessExpireTime  time.Duration `mapstructure:"ACCESS_EXPIRE_TIME"`
		RefreshExpireTime time.Duration `mapstructure:"REFRESH_EXPIRE_TIME"`
	} `mapstructure:"JWT"`

	Database struct {
		RelaySQL DatabaseConfig `mapstructure:"RELAY_SQL"`
	} `mapstructure:"DATABASE"`

	HTMLTemplate struct {
		SystemMaintenance string `mapstructure:"SYSTEM_MAINTENANCE"`
	} `mapstructure:"HTML_TEMPLATE"`
}

// InitConfig init config
func InitConfig() error {
	v := viper.New()
	v.AddConfigPath(filePath)
	v.SetConfigName(fileNameConfig)
	v.SetConfigType(fileExtension)
	v.AutomaticEnv()

	// แปลง _ underscore เป็น . dot
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		logger.Log.Errorf("read config file error: %s", err)
		return err
	}

	if err := bindingConfig(v, CF); err != nil {
		logger.Log.Errorf("binding config error: %s", err)
		return err
	}

	// set config ปิด/เปิด ระบบ
	if err := initConfigAvailable(); err != nil {
		logger.Log.Errorf("init config available error: %s", err)
		return err
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		logger.Log.Infof("config file changed: %s", e.Name)
		if err := v.Unmarshal(CF); err != nil {
			logger.Log.Errorf("binding config error: %s", err)
		}
	})
	v.WatchConfig()

	return nil
}

// bindingConfig binding config
func bindingConfig(vp *viper.Viper, cf *Configs) error {
	if err := vp.Unmarshal(&cf); err != nil {
		logger.Log.Errorf("unmarshal config error: %s", err)
		return err
	}

	validate := validator.New()
	if err := validate.RegisterValidation("maxString", validateString); err != nil {
		logger.Log.Errorf("cannot register maxString Validator config error: %s", err)
		return err
	}

	en := en.New()
	cf.UniversalTranslator = ut.New(en, en)
	enTrans, _ := cf.UniversalTranslator.GetTranslator("en")
	if err := en_translations.RegisterDefaultTranslations(validate, enTrans); err != nil {
		logger.Log.Errorf("cannot add english translator config error: %s", err)
		return err
	}

	_ = validate.RegisterTranslation("maxString", enTrans, func(ut ut.Translator) error {
		return ut.Add("maxString", "Sorry, {0} cannot exceed {1} characters", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		field := strings.ToLower(fe.Field())
		t, _ := ut.T("maxString", field, fe.Param())
		return t
	})

	cf.Validator = validate

	return nil
}

// validateString implements validator.Func for max string by rune
func validateString(fl validator.FieldLevel) bool {
	var err error
	limit := 255
	param := strings.Split(fl.Param(), ":")
	if len(param) > 0 {
		limit, err = strconv.Atoi(param[0])
		if err != nil {
			limit = 255
		}
	}

	if lengthOfString := utf8.RuneCountInString(fl.Field().String()); lengthOfString > limit {
		return false
	}

	return true
}

// initConfigAvailable init config available
// init config ปิด/เปิด ระบบ
func initConfigAvailable() error {
	// create file config
	if err := CF.SetConfigAvailableStatus(AvailableStatusOnline); err != nil {
		logger.Log.Errorf("creating file available status error: %s", err)
		return err
	}

	if err := CF.SetConfigAvailableDescription(""); err != nil {
		logger.Log.Errorf("creating file available description error: %s", err)
		return err
	}

	// read config
	v := viper.New()
	v.AddConfigPath(filePath)
	v.SetConfigName(fileNameConfigAvailableStatus)
	v.SetConfigType(fileExtension)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logger.Log.Errorf("read config file error: %s", err)
		return err
	}

	cf := &AvailableConfig{}
	if err := v.Unmarshal(cf); err != nil {
		logger.Log.Errorf("binding config error: %s", err)
		return err
	}
	CF.App.AvailableStatus = cf.Status

	v.OnConfigChange(func(e fsnotify.Event) {
		logger.Log.Infof("config file changed: %s", e.Name)
		if err := v.Unmarshal(cf); err != nil {
			logger.Log.Errorf("binding config error: %s", err)
		}
		CF.App.AvailableStatus = cf.Status
	})
	v.WatchConfig()

	return nil
}

// SetConfigAvailableStatus set config available status
// สร้าง config สถานะ ปิด/เปิด ระบบ
func (cf *Configs) SetConfigAvailableStatus(status string) error {
	d, _ := sonic.Marshal(AvailableConfig{
		Status: status,
	})
	p := fmt.Sprintf("%s/%s", filePath, fileNameConfigAvailableStatus)
	err := os.WriteFile(p, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

// SetConfigAvailableDescription set config available description
// สร้าง config html ใช้แสดงเมื่อปิดระบบ
func (cf *Configs) SetConfigAvailableDescription(body string) error {
	d := []byte(body)
	p := fmt.Sprintf("%s/%s", filePath, fileNameConfigAvailableDescription)
	err := os.WriteFile(p, d, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ReadConfigAvailableDescription read config available description
// อ่าน config html ใช้แสดงเมื่อปิดระบบ
func (cf *Configs) ReadConfigAvailableDescription() (string, error) {
	d, err := os.ReadFile(fmt.Sprintf("./%s/%s", filePath, fileNameConfigAvailableDescription))
	if err != nil {
		return "", err
	}

	return string(d), nil
}
