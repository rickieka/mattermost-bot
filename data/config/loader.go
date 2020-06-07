package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Conf contains all values fron conf/conf.yml
var Conf Config

// InitConfig reads the conf/conf.yml file and initializes the Conf variable
func InitConfig(configFile string) {
	readConfigFile(configFile)
}

// Config Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type Config struct {
	Mattermost struct {
		HttpUrl      string `yaml:"host"`
		WebsocketUrl string `yaml:"websocket"`
		Channels     struct {
			Debugging string `yaml:"debugging"`
			Team      string `yaml:"team"`
		}
	}
	Bot struct {
		ID          string `yaml:"id"`
		AccessToken string `yaml:"access_token"`
		Name        string `yaml:"name"`
	}
	Mysql   []SqlConfig
	Logging struct {
		Enabled bool   `yaml:enabled`
		Logfile string `yaml:logfile`
	}
}

type SqlConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:database`
}

// Read config.yaml
// If the config can not be read correctly, exit the program
func readConfigFile(configFile string) {
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)

	}
	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
