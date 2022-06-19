package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	keyJira               = "jira"
	keyBaseUrl            = "base_url"
	keyToken              = "token"
	keyProjects           = "projects"
	keyDays               = "days"
	keyMail               = "smtp"
	keyMailSrv            = "address"
	keyMailSender         = "account"
	keyMailSenderPassword = "password"
	keyMailReceiver       = "to"
)

type Config struct {
}

func GetConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	return Config{}, nil
}

func (c *Config) GetBaseURL() string {
	return viper.GetStringMapString(keyJira)[keyBaseUrl]
}

func (c *Config) GetAccessToken() string {
	return viper.GetStringMapString(keyJira)[keyToken]
}

func (c *Config) GetProjects() (map[string]map[string][]string, error) {
	res := map[string]map[string][]string{}

	projects := viper.GetStringSlice(keyProjects)
	for _, project := range projects {
		res[project] = map[string][]string{}

		fields := viper.GetStringMapStringSlice(project)
		for key, value := range fields {
			res[project][key] = value
		}
	}

	return res, nil
}

func (c *Config) GetDays() uint {
	return viper.GetUint(keyDays)
}

func (c *Config) GetMailServer() string {
	return viper.GetStringMapString(keyMail)[keyMailSrv]
}

func (c *Config) GetMailSender() string {
	return viper.GetStringMapString(keyMail)[keyMailSender]
}

func (c *Config) GetMailSenderPassword() string {
	return viper.GetStringMapString(keyMail)[keyMailSenderPassword]
}

func (c *Config) GetMailReciever() []string {
	return viper.GetStringMapStringSlice(keyMail)["to"]
}
