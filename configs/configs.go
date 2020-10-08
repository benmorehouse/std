// Package configs provides basic, configurable data needed to run std
package configs

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	userOS "os/user"
	"path/filepath"
)

type config struct {
	DatabasePath   string `yaml:"database_path"`
	TempWorkSpace  string `yaml:"temp_workspace"`
	LogFile        string `yaml:"log_file"`
	BufferMDFile   string `yaml:"buffer_md_file"`
	VaultAddr      string `yaml:"vault_addr"`
	VaultStaticKey string `yaml:"vault_static_key"`
	dotFileConfig  string
}

const (
	// ListBucketLabel is the bucket label which holds all the secret sauce
	ListBucketLabel = "std_bucket_list"
	// ListBucketKey is the key within the list bucket label
	ListBucketKey = "std_bucket_key"
)

// STDConf embodies the set configurations for the application
var STDConf = &config{
	DatabasePath:   ".std/db/mainDB.db",
	TempWorkSpace:  ".std/tmp/",
	LogFile:        ".std/logging/std.log",
	BufferMDFile:   ".std/logging/stdin.buffer.md",
	VaultAddr:      "",
	VaultStaticKey: "",
	dotFileConfig:  ".std/config.yml",
}

var root string

func init() {
	if err := setConfigWithUserRoot(); err != nil {
		log.Fatal("couldn't set the necessary directory names", err)
	}

	if err := makePaths(); err != nil {
		log.Fatal("couldn't make the neccessary paths", err)
	}

	if err := loadDotfileConfig(); err != nil {
		log.Fatal("Cannot load yaml config:", err)
	}
}

func setConfigWithUserRoot() error {
	usr, err := userOS.Current()
	if err != nil {
		return err
	}

	root = usr.HomeDir
	STDConf.dotFileConfig = filepath.Join(root, STDConf.dotFileConfig)
	STDConf.BufferMDFile = filepath.Join(root, STDConf.BufferMDFile)
	STDConf.DatabasePath = filepath.Join(root, STDConf.DatabasePath)
	STDConf.LogFile = filepath.Join(root, STDConf.LogFile)
	STDConf.TempWorkSpace = filepath.Join(root, STDConf.TempWorkSpace)
	return nil
}

func makePaths() error {
	if err := os.MkdirAll(filepath.Dir(STDConf.dotFileConfig), 0744); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(STDConf.BufferMDFile), 0744); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(STDConf.DatabasePath), 0744); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(STDConf.LogFile), 0744); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(STDConf.TempWorkSpace), 0744); err != nil {
		return err
	}
	return nil
}

func loadDotfileConfig() error {
	rawBytes, err := ioutil.ReadFile(STDConf.dotFileConfig)
	if err != nil {
		if err := writeDefaultConfig(STDConf.dotFileConfig); err != nil {
			return fmt.Errorf("write_default_config_fail: %s", err.Error())
		}
		return nil
	}

	if err := yaml.Unmarshal(rawBytes, &STDConf); err != nil {
		return fmt.Errorf("yaml_config_unmarshal: %s", err.Error())
	}
	return nil
}

// writeDefaultConfig will write a config in the default location
func writeDefaultConfig(configFileName string) error {
	if err := os.MkdirAll(filepath.Dir(configFileName), 0744); err != nil {
		return fmt.Errorf("make_all_default_config_fail: %s", err.Error())
	}

	rawBytes, err := yaml.Marshal(STDConf)
	if err != nil {
		return fmt.Errorf("default_config_marshal_fail: %s", err.Error())
	}

	if err := ioutil.WriteFile(configFileName, rawBytes, 0744); err != nil {
		return fmt.Errorf("default_config_marshal_write_fail: %s", err.Error())
	}

	return nil
}
