package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultEnv  = "local"
	defaultPort = 8080
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	OSS      OSSConfig      `mapstructure:"oss"`
	Security SecurityConfig `mapstructure:"security"`
	Log      LogConfig      `mapstructure:"log"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
}

type PostgresConfig struct {
	Host                   string `mapstructure:"host"`
	Port                   int    `mapstructure:"port"`
	User                   string `mapstructure:"user"`
	Password               string `mapstructure:"password"`
	DBName                 string `mapstructure:"dbname"`
	SSLMode                string `mapstructure:"sslmode"`
	MaxOpenConns           int    `mapstructure:"max_open_conns"`
	MaxIdleConns           int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeMinutes int    `mapstructure:"conn_max_lifetime_minutes"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type SecurityConfig struct {
	Behavior BehaviorSecurityConfig `mapstructure:"behavior"`
}

type BehaviorSecurityConfig struct {
	Enabled                    bool  `mapstructure:"enabled"`
	WindowSeconds              int64 `mapstructure:"window_seconds"`
	IPLimitPerWindow           int   `mapstructure:"ip_limit_per_window"`
	SuspiciousIPLimitPerWindow int   `mapstructure:"suspicious_ip_limit_per_window"`
}

type OSSConfig struct {
	BucketName          string `mapstructure:"bucket_name"`
	Endpoint            string `mapstructure:"endpoint"`
	PublicBaseURL       string `mapstructure:"public_base_url"`
	Region              string `mapstructure:"region"`
	PresignExpireSecond int    `mapstructure:"presign_expire_seconds"`
}

func Load() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = defaultEnv
	}

	v := viper.New()
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	if cfg.App.Env == "" {
		cfg.App.Env = env
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "luke-chu-site-api")
	v.SetDefault("app.env", defaultEnv)
	v.SetDefault("app.port", defaultPort)
	v.SetDefault("server.read_timeout", 10)
	v.SetDefault("server.write_timeout", 15)
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.max_open_conns", 20)
	v.SetDefault("postgres.max_idle_conns", 10)
	v.SetDefault("postgres.conn_max_lifetime_minutes", 30)
	v.SetDefault("oss.presign_expire_seconds", 300)
	v.SetDefault("security.behavior.enabled", true)
	v.SetDefault("security.behavior.window_seconds", 60)
	v.SetDefault("security.behavior.ip_limit_per_window", 120)
	v.SetDefault("security.behavior.suspicious_ip_limit_per_window", 20)
	v.SetDefault("log.level", "info")
}
