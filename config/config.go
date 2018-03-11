package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var envKeys = []string{
	"PG_HOST", "PG_PORT", "PG_USER", "PG_PASS", "PG_NAME",
	"CS_CLUSTER", "CS_PORT", "CS_USER", "CS_PASSWORD", "CS_KEYSPACE",
	"RD_HOST", "RD_PORT", "RD_PASSWORD",
	"GB_PORT",
}

// Config app configuration
type Config struct {
	Postgresql Postgresql
	Cassandra  Cassandra
	Redis      Redis
	App        App
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
	Cluster  string `mapstructure:"cs_cluster"`
	Port     string `mapstructure:"cs_port"`
	User     string `mapstructure:"cs_user"`
	Password string `mapstructure:"cs_password"`
	Keyspace string `mapstructure:"cs_keyspace"`
}

type Redis struct {
	Host     string `mapstructure:"rd_host"`
	Port     string `mapstructure:"rd_port"`
	Password string `mapstructure:"rd_password"`
}

// App
type App struct {
	Port int `mapstructure:"gb_port"`
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
		fmt.Printf("%+v", cf)

		return &cf, nil
	}
	viper.AutomaticEnv()
	for _, key := range envKeys {
		viper.BindEnv(key)
	}
	if err := viper.Unmarshal(&cf.Postgresql); err != nil {
		return nil, err
	}

	for _, key := range envKeys {
		viper.BindEnv(key)
	}
	if err := viper.Unmarshal(&cf.Redis); err != nil {
		return nil, err
	}

	for _, key := range envKeys {
		viper.BindEnv(key)
	}
	if err := viper.Unmarshal(&cf.Cassandra); err != nil {
		return nil, err
	}

	for _, key := range envKeys {
		viper.BindEnv(key)
	}
	if err := viper.Unmarshal(&cf.App); err != nil {
		return nil, err
	}

	return &cf, nil
}
