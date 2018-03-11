package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var envKeys = []string{"PG_HOST", "PG_PORT", "PG_USER", "PG_PASS", "PG_NAME", "APP_PORT"}

// Config app configuration
type Config struct {
	Postgresql Postgresql `mapstructure:"postgresql"`
	Cassandra  Cassandra  `mapstructure:"cassandra"`
	Redis      Redis      `mapstructure:"redis"`
	App        App        `mapstructure:"app"`
}

// Postgresql Postgresql configuration
type Postgresql struct {
	Host     string `mapstructure:"pg_host"`
	Port     string `mapstructure:"pg_port"`
	User     string `mapstructure:"pg_user"`
	Password string `mapstructure:"pg_pass"`
	Name     string `mapstructure:"pg_name"`
}

type Cassandra struct {
	Cluster  string `mapstructure:"cluster"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Keyspace string `mapstructure:"keyspace"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

// App
type App struct {
	Port int `mapstructure:"port"`
}

// ReadConfig read configuration from ENV or from config file
func ReadConfig(configFile string) (*Config, error) {
	var cf Config
	if configFile != "" {
		// read from config yaml
		viper.AddConfigPath(".")
		viper.SetConfigName(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
		if err := viper.Unmarshal(&cf); err != nil {
			panic(err)
		}

		return &cf, nil
	}
	viper.SetEnvPrefix("GB")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	for _, key := range envKeys {
		fmt.Println(key, viper.Get(key))
		viper.BindEnv(key)
	}
	if err := viper.Unmarshal(&cf); err != nil {
		return nil, err
	}

	return &cf, nil
}
