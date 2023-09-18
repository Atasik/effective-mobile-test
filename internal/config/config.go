package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort                = "8080"
	defaultHTTPRWTimeout           = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes  = 1
	defaultDatabaseRefreshInterval = 30 * time.Second
)

type (
	Config struct {
		Postgres       PostgresConfig
		HTTP           HTTPConfig
		Kafka          KafkaConfig
		Redis          RedisConfig
		CacheTTL       time.Duration `mapstructure:"ttl"`
		KafkaEndpoints []string      `mapstructure:"kafka-endpoints"`
		KafkaTopics    []string      `mapstructure:"kafka-topics"`
	}

	PostgresConfig struct {
		Username string
		Password string
		Port     string
		Host     string
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	HTTPConfig struct {
		Host               string
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegaBytes int           `mapstructure:"maxHeaderMegaBytes"`
	}

	RedisConfig struct {
		Address  string
		Password string
		DB       int `mapstructure:"db"`
	}

	KafkaConfig struct {
		GroupID        string `mapstructure:"group-id"`
		ClientID       string `mapstructure:"client-id"`
		TLSEnable      bool   `mapstructure:"tls-enable"`
		ReturnSucceses bool   `mapstructure:"returnSucceses"`
		RequiredAcks   int    `mapstructure:"requiredAcks"`
	}
)

func InitConfig(configPath string) (*Config, error) {
	setDefaults()

	if err := parseConfigFile(configPath); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("cache.ttl", &cfg.CacheTTL); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("kafka", &cfg.Kafka); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("kafka-endpoints", &cfg.KafkaEndpoints); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("kafka-topics", &cfg.KafkaTopics); err != nil {
		return err
	}

	return viper.UnmarshalKey("http", &cfg.HTTP)
}

func setFromEnv(cfg *Config) {
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Redis.Address = os.Getenv("REDIS_ADDRESS")
	cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}

func setDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderMegaBytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.readTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("postgres.refreshInterval", defaultDatabaseRefreshInterval)
}
