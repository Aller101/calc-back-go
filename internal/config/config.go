package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
	PostgresDB `yaml:"postgres_db"`
}

type HTTPServer struct {
	Address string `yaml:"address"`
}

type PostgresDB struct {
	User     string `yaml:"user"`
	Dbname   string `yaml:"dbname"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Sslmode  string `yaml:"sslmode"`
}

func MustLoad() *Config {
	os.Setenv("CONFIG_PATH", "./configs/local.yaml") //пдлхо
	conf_path := os.Getenv("CONFIG_PATH")
	if conf_path == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(conf_path); err != nil {
		log.Fatalf("config file %s does not exist", conf_path)
	}

	var conf Config
	err := cleanenv.ReadConfig(conf_path, &conf)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &conf
}
