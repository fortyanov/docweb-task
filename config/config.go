package config

import (
	"os"
	"path"

	"github.com/go-ini/ini"

	"docweb-task/server"
	"docweb-task/storage"
)

type Config struct {
	General General        `ini:"general"`
	Server  server.Config  `ini:"server"`
	Storage storage.Config `ini:"storage"`
}

type General struct {
	LogLevel string `ini:"log_level"`
	LogDest  string `ini:"log_destination"`
	LogTag   string `ini:"log_tag"`
	PidFile  string `ini:"pidfile"`
}

func (cfg *Config) Validate() (err error) {
	if err = cfg.Server.Validate(); err != nil {
		return err
	}
	if err = cfg.Storage.Validate(); err != nil {
		return err
	}

	return err
}

func Parse(fileName string) (*Config, error) {
	var (
		err     error
		iniFile *ini.File
	)

	if iniFile, err = ini.Load(fileName); err != nil {
		return nil, err
	}

	cfg := &Config{
		General: General{
			LogLevel: "debug",
			LogTag:   path.Base(os.Args[0]),
			PidFile:  "/var/run/" + path.Base(os.Args[0]) + ".pid",
		},
	}

	if err = iniFile.MapTo(cfg); err != nil {
		return nil, err
	}

	return cfg, cfg.Validate()
}
