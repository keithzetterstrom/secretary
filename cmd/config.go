package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"github.com/keithzetterstrom/secretary/internal/repository/docs"
	"github.com/keithzetterstrom/secretary/internal/tgbot"
	"github.com/keithzetterstrom/secretary/utils/logger"
)

type Config struct {
	ServiceConfig `yaml:"service"`
	BotConfig     tgbot.Config `yaml:"bot"`
	DocsConfig    docs.Config  `yaml:"docs"`
}

type ServiceConfig struct {
	Env  string
	Host string            `yaml:"host"`
	Port string            `yaml:"port"`
	Logs logger.LogsConfig `yaml:"logs"`
}

var (
	env     = flag.String("env", "prod", "")
	cnfPath = flag.String("config", "", "")
	netAddr = flag.String("addr", "", "")
	envFile = flag.String("env_file", "", "")
)

func NewConfig(cfg *Config) error {
	flag.Parse()

	if *cnfPath == "" {
		return errors.New("no config path")
	}

	f, err := os.Open(*cnfPath)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	cfg.Env = *env

	if *envFile != "" {
		err := godotenv.Load(*envFile)
		if err != nil {
			return err
		}
	}

	if *netAddr != "" {
		cfg.Host = *netAddr
	}

	cfg.BotConfig.Token = os.Getenv(cfg.BotConfig.Token)

	cfg.DocsConfig.SpreadsheetId = os.Getenv(cfg.DocsConfig.SpreadsheetId)

	return nil
}
