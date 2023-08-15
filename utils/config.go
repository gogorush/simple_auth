package utils

import (
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

var globalConfig *FileConfig
var onceConfig sync.Once

func GetConfig() *FileConfig {
	return globalConfig
}

func ReadConfig(configFile string, conf interface{}) error {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		panic(err)
	}
	return nil
}

// Config : global config struct
type FileConfig struct {
	Address string `yaml:"Address,omitempty"`
	JwtKey  string `yaml:"JwtKey,omitempty"`
	Admin   Admin  `yaml:"Admin,omitempty"`
}

type Admin struct {
	UserName string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// InitConfig : get a new Config struct.
func InitConfig(confFile string) error {
	if err := ReadConfig(confFile, &globalConfig); err != nil {
		return err
	}
	return nil
}
