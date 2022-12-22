package config

import (
	"fmt"
	"path/filepath"

	log "github.com/ahmadmirdas/julo-test/utils/log"
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Environment is the running application mode (dev, test, production)
var Environment string = "dev"
var Config config

type config struct {
	Environment string `mapstructure:"environment"`
	PostgresCfg struct {
		Database    string `mapstructure:"database"`
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		Username    string `mapstructure:"username"`
		Password    string `mapstructure:"password"`
		MaxConn     int    `mapstructure:"max_conn"`
		MinIdleConn int    `mapstructure:"min_idle_conn"`
		MaxRetries  int    `mapstructure:"max_retries"`
	} `mapstructure:"postgres"`
	JWTCfg struct {
		Issuer  string `mapstructure:"issuer"`
		Exp     int    `mapstructure:"exp"`
		SignKey string `mapstructure:"sign_key"`
	} `mapstructure:"jwt"`
}

func init() {
	var err error

	configureLogging()

	viper.SetEnvPrefix("JT")
	viper.AutomaticEnv()

	configName := "application"

	if viper.IsSet("ENV") {
		Environment = viper.GetString("ENV")
	}

	if Environment != "production" {
		configName = configName + "." + Environment
	}

	logrus.Info("JT_ENV: ", Environment)

	// name of config file (without extension)
	viper.SetConfigName(configName)
	viper.AddConfigPath(filepath.Join(GetAppBasePath(), "config/app"))
	// optionally look for config in the working directory
	viper.AddConfigPath(".")

	// Find and read the config file
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	//Unmarshal application yml to config
	err = viper.Unmarshal(&Config)
	if err != nil {
		logrus.Errorf("unable to decode into struct, %v", err)
	}
}

func GetAppBasePath() string {
	basePath, _ := filepath.Abs(".")
	for filepath.Base(basePath) != "julo-test" {
		basePath = filepath.Dir(basePath)
	}

	return basePath
}

func configureLogging() {
	logrus.AddHook(log.LogrusSourceContextHook{})
	logrus.SetFormatter(&joonix.Formatter{})
	logrus.SetLevel(logrus.DebugLevel)
}
