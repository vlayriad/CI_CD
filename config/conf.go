package config

import (
	"encoding/json"
	"fmt"
	"kaffein/utils/logger"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Config struct {
	App struct {
		AppName        string `json:"app_name"`
		AppVersion     string `json:"app_version"`
		AppEnvironment string `json:"app_environment"`
		AppDebugMode   string `json:"app_debug_mode"`
		AppTimeZone    string `json:"app_time_zone"`
		AppHost        string `json:"app_host"`
		AppPort        int    `json:"app_port"`
		AppJWT         string `json:"app_jwt"`
	} `json:"server"`

	Postgres struct {
		PGHost    string `json:"host"`
		PGPort    int    `json:"port"`
		PGUser    string `json:"user"`
		PGPass    string `json:"password"`
		PGDB      string `json:"dbname"`
		PGSSLMode bool   `json:"sslmode"`
	} `json:"postgres"`

	Redis struct {
		RedisHost     string `json:"host"`
		RedisPort     int    `json:"port"`
		RedisPassword string `json:"password"`
		RedisDB       int    `json:"db"`
	} `json:"redis"`

	Minio struct {
		Endpoint        string `json:"endpoint"`
		BucketName      string `json:"bucket_name"`
		AccessKeyID     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
		UseSSL          bool   `json:"use_ssl"`
		OriginEndpoint  string `json:"origin_endpoint"`
	} `json:"minio"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	if instance == nil {
		log.Panic().
			Msg("config not initialized — call SetupApp() first")
	}
	return instance
}

func SetupApp(fileName string) error {
	var err error
	once.Do(func() {
		logger.InitDefault()

		var cfg Config
		content, e := os.ReadFile(fileName)
		if e != nil {
			log.Error().
				Err(e).
				Msg("failed to read config file")
			err = fmt.Errorf("failed to read config file: %w", e)
			return
		}

		if e := json.Unmarshal(content, &cfg); e != nil {
			log.Error().
				Err(e).
				Msg("failed to parse config file")
			err = fmt.Errorf("failed to parse config file: %w", e)
			return
		}

		instance = &cfg

		loc, e := time.LoadLocation(cfg.App.AppTimeZone)
		if e != nil {
			log.Warn().
				Err(e).
				Msgf("invalid timezone '%s', fallback to Asia/Jakarta", cfg.App.AppTimeZone)
			loc, _ = time.LoadLocation("Asia/Jakarta")
		}
		time.Local = loc
		logger.Init(cfg.App.AppTimeZone, cfg.App.AppEnvironment, cfg.App.AppDebugMode)

		log.Info().
			Str("timezone", time.Now().Location().String()).
			Str("mode", cfg.App.AppDebugMode).
			Msg("Logger and system timezone initialized")
	})
	return err
}
