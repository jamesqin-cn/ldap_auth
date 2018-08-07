package utils

import (
	"io/ioutil"
	"log"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type ServerConf struct {
	Listen string `yaml:"listen"`
}

type GlobalConfig struct {
	LdapConf   LdapConf   `yaml:"ldap"`
	ServerConf ServerConf `yaml:"http"`
}

var (
	cfgCache *GlobalConfig
)

//GetConfig get cfg
func GetConfig(cfgFile string) *GlobalConfig {
	if cfgCache != nil {
		return cfgCache
	}

	cfg := &GlobalConfig{}

	fileFullName, _ := filepath.Abs(cfgFile)
	yamlFile, err := ioutil.ReadFile(fileFullName)
	log.Println("load config ", fileFullName)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		panic(err)
	}

	cfgCache = cfg

	return cfg
}
