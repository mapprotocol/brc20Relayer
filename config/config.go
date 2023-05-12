// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

const configFile = "/config.json"

func DefaultConfigFile() string {
	return DefaultDir() + configFile
}

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) <= 0
}

type Config struct {
	URL            string         `json:"url"`
	URL2           string         `json:"url2"`
	Contract       string         `json:"contract"`
	Start          uint64         `json:"start"`
	Count          uint64         `json:"count"`
	DatabaseConfig DataBaseConfig `json:"databaseConfig"`
}

type DataBaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c *Config) validate() error {
	if IsEmpty(c.URL) {
		return fmt.Errorf("required field URL")
	}
	if IsEmpty(c.URL2) {
		return fmt.Errorf("required field URL2")
	}
	if IsEmpty(c.Contract) {
		return fmt.Errorf("required field contract")
	}
	// todo validate database config

	return nil
}

func GetConfig(ctx *cli.Context) (*Config, error) {
	var cfg Config
	path := DefaultConfigFile()
	if file := ctx.String(ConfigFileFlag.Name); file != "" {
		path = file
	}
	err := loadConfig(path, &cfg)
	if err != nil {
		log.Warn("failed to loading json file", "err", err.Error())
		return &cfg, err
	}

	log.Debug("Loaded config", "path", path)
	err = cfg.validate()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadConfig(file string, config *Config) error {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	log.Debug("Loading configuration", "path", filepath.Clean(fp))

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}

	if ext == ".json" {
		if err = json.NewDecoder(f).Decode(&config); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}

	return nil
}
