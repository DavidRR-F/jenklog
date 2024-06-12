package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sync"
)

const (
	CONFIG_FILE_NAME = ".env"
	CONFIG_FILE_TYPE = "env"
)

type Config struct {
	Username string `mapstructure:"JENKINS_USER_ID"`
	Token    string `mapstructure:"JENKINS_API_TOKEN"`
	URL      string `mapstructure:"JENKINS_URL"`
}

var (
	instance *Config
	once     sync.Once
	err      error
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		instance = &Config{}
		err = instance.loadConfig()
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (c *Config) loadConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(homeDir, CONFIG_FILE_NAME)

	viper.SetConfigFile(configFilePath)
	viper.SetConfigType(CONFIG_FILE_TYPE)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

func SaveConfig(username, token, url string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(homeDir, CONFIG_FILE_NAME)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		file, err := os.Create(configFilePath)
		if err != nil {
			return err
		}
		file.Close()
	}

	viper.Set("JENKINS_USER_ID", username)
	viper.Set("JENKINS_API_TOKEN", token)
	viper.Set("JENKINS_URL", url)

	viper.SetConfigFile(configFilePath)
	viper.SetConfigType(CONFIG_FILE_TYPE)

	return viper.WriteConfigAs(configFilePath)
}
