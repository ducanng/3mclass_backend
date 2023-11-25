package config

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/spf13/viper"
)

//go:embed config.yaml
var defaultConfig []byte

type (
	Config struct {
		HTTPPort          string            `mapstructure:"http_port"`
		DSN               string            `mapstructure:"db_connection_string"`
		CORS              cors              `mapstructure:"cors"`
		GormDebug         bool              `mapstructure:"gorm_debug"`
		DBAutoMigrate     bool              `mapstructure:"db_auto_migrate"`
		JWT               jwt               `mapstructure:"jwt"`
		AccessTokenCookie accessTokenCookie `mapstructure:"access_token_cookie"`
		SecretKey         string            `mapstructure:"secret_key"`
		BaseHost          string            `mapstructure:"base_host"`
	}
	cors struct {
		Enabled          bool     `mapstructure:"enabled"`
		Origins          []string `mapstructure:"origins"`
		AllowedMethods   []string `mapstructure:"allowed_methods"`
		ExposedHeaders   []string `mapstructure:"exposed_headers"`
		AllowedHeaders   []string `mapstructure:"allowed_headers"`
		Debug            bool     `mapstructure:"debug"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
	}
	jwt struct {
		ExpiryTime int64  `mapstructure:"expiry_time"`
		PrivateKey string `mapstructure:"private_key"`
		PublicKey  string `mapstructure:"public_key"`
		Enabled    bool   `mapstructure:"enabled"`
	}
	accessTokenCookie struct {
		Domain         string `mapstructure:"domain"`
		CookieName     string `mapstructure:"cookie_name"`
		PreviousDomain string `mapstructure:"previous_domain"`
	}
)

func Load() *Config {
	cfg := &Config{}
	v := viper.New()
	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(defaultConfig))
	if err != nil {
		log.Panicf("Failed to read viper config, err=%s", err.Error())
	}
	if err = v.Unmarshal(&cfg); err != nil {
		log.Panicf("Failed to unmarshal config, err=%s", err.Error())
	}
	return cfg
}
