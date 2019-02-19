package config

import (
	"io/ioutil"
	"strconv"
	"strings"
)

var config *Config = nil

const cfgPath = "/etc/runner.cfg"

type Config struct {
	FilesPath    string
	DBName       string
	DBUser       string
	DBPassword   string
	WhereToStore string
	ComiledPath  string
	LogLevel     int
	LogType      string
	LogPath      string
}

func InitConfig() (*Config, error) {
	bytes, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	strings := strings.Split(string(bytes), "\n")

	level, err := strconv.Atoi(strings[6])

	return &Config{
		strings[0],
		strings[1],
		strings[2],
		strings[3],
		strings[4],
		strings[5],
		level,
		strings[7],
		strings[8],
	}, nil

}

func (conf Config) GetLoggerType() string {
	return conf.LogType
}

func (conf Config) GetLoggerLevel() int {
	return conf.LogLevel
}

func (conf Config) GetLoggerPath() string {
	return conf.LogPath
}

func GetConfigInstance() (*Config, error) {
	if config == nil {
		return InitConfig()
	}
	return config, nil
}
