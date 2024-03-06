package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ApplicationConfig ApplicationConfig `mapstructure:"APP_CONFIG"`
	DBFirstConfig     DBConfig          `mapstructure:"DB_CONFIG_1"`
	DBSecondConfig    DBConfig          `mapstructure:"DB_CONFIG_2"`
}

type ApplicationConfig struct {
	Host string `mapstructure:"APP_HOST"`
	Port string `mapstructure:"APP_PORT"`
}

type DBConfig struct {
	Driver             string        `mapstructure:"DRIVER"`
	Username           string        `mapstructure:"USERNAME"`
	Password           string        `mapstructure:"PASSWORD"`
	Name               string        `mapstructure:"NAME"`
	Host               string        `mapstructure:"HOST"`
	Port               int           `mapstructure:"PORT"`
	MaxIdleConnections int           `mapstructure:"MAX_IDLE_CONNECTIONS"`
	MaxOpenConnections int           `mapstructure:"MAX_OPEN_CONNECTIONS"`
	MaxConnLifetime    time.Duration `mapstructure:"MAX_CONN_LIFETIME"`
	DebugMode          bool          `mapstructure:"DEBUG_MODE"`
	Timeout            string        `mapstructure:"TIMEOUT"`
	WriteTimeout       string        `mapstructure:"WRITE_TIMEOUT"`
	ReadTimeout        string        `mapstructure:"READ_TIMEOUT"`
	SSLMode            string        `mapstructure:"SSLMODE"`
}

func Load() (conf Config) {
	viper.SetConfigFile("config")
	viper.SetConfigFile("./config.json")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}

	return
}

func (db *DBConfig) GetDSN() (dsn string) {
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&writeTimeout=%s&readTimeout=%s&charset=utf8mb4&parseTime=True&loc=Local",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.Timeout,
		db.WriteTimeout,
		db.ReadTimeout,
	)

	if db.Driver == "postgres" {
		dsn = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
			db.Driver,
			db.Username,
			db.Password,
			db.Host,
			db.Port,
			db.Name,
			db.SSLMode,
		)
	}

	fmt.Println(dsn)

	return
}
