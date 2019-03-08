package config

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	iot "io/ioutil"
	"strings"
)

const (
	HTTP_USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.90 Safari/537.36"
	HTTP_ACCEPT     = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
)

var UrlConfig *Config

func init() {
	var err error

	if UrlConfig == nil {
		UrlConfig, err = LoadConf("conf/app.yml")
		if err != nil {
			panic(err)
		}
	}

}

func CreateDatabase() (*sql.DB, error) {
	//UrlConfig := CreateConfig()
	//fmt.Println("UrlConfig", UrlConfig)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		UrlConfig.Postgres.Host, UrlConfig.Postgres.Port, UrlConfig.Postgres.User, UrlConfig.Postgres.Password, UrlConfig.Postgres.DB)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db, nil
}

// Config contains the configuration of the url shortener.
type Config struct {
	Httpport   string `yaml:"http_port" json:"http_port"`
	Grpclisten string `yaml:"grpc_listen" json:"grpc_listen"`
	Server     struct {
		Host string `yaml:"host" json:"host"`
		Port string `yaml:"port" json:"port"`
	} `yaml:"server" json:"server"`
	Postgres struct {
		Host     string `yaml:"host" json:"host"`
		Port     string `yaml:"port" json:"port"`
		User     string `yaml:"user" json:"user"`
		Password string `yaml:"password" json:"password"`
		DB       string `yaml:"db" json:"db"`
	} `yaml:"postgres" json:"postgres"`
	Options struct {
		Prefix string `yaml:"prefix" json:"prefix"`
	} `yaml:"options" json:"options"`
}

func LoadConf(filepath string) (*Config, error) {
	if filepath == "" {
		return nil, errors.New("filepath is empty, must use --config xxx.yml/json")
	}

	data, err := iot.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if strings.HasSuffix(filepath, ".json") {
		err = json.Unmarshal(data, &cfg)
	} else if strings.HasSuffix(filepath, ".yml") || strings.HasSuffix(filepath, ".yaml") {
		err = yaml.Unmarshal(data, &cfg)
	} else {
		return nil, errors.New("you config file must be json/yml")
	}

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
