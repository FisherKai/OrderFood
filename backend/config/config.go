package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	JWT          JWTConfig          `mapstructure:"jwt"`
	Storage      StorageConfig      `mapstructure:"storage"`
	Announcement AnnouncementConfig `mapstructure:"announcement"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type StorageConfig struct {
	LocalPath       string   `mapstructure:"local_path"`
	ImageMaxSize    int64    `mapstructure:"image_max_size"`
	AllowedFormats  []string `mapstructure:"allowed_formats"`
	ThumbnailWidth  int      `mapstructure:"thumbnail_width"`
	ThumbnailHeight int      `mapstructure:"thumbnail_height"`
}

type AnnouncementConfig struct {
	AutoPlayInterval int `mapstructure:"auto_play_interval"`
	MaxDisplayCount  int `mapstructure:"max_display_count"`
}

var AppConfig *Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	AppConfig = &Config{}
	return viper.Unmarshal(AppConfig)
}
