// Package configs provides basic, configurable data needed to run std
package configs

import (
	"fmt"
	"io/ioutil"
	"os"
	userOS "os/user"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	ListDatabasePath     string `yaml:"list_database_path"`
	PasswordDatabasePath string `yaml:"password_database_path"`
	TempWorkSpace        string `yaml:"temp_workspace"`
	LogFile              string `yaml:"log_file"`
	BufferMDFile         string `yaml:"buffer_md_file"`
	dotFileConfig        string
}

const (
	// ListBucketLabel is the bucket label which holds all the secret sauce
	ListBucketLabel = "std_bucket_list"
	// ListBucketKey is the key within the list bucket label
	ListBucketKey = "std_bucket_key"
)

const (
	testingConfigEnv            = "STD_CONFIG_TESTING"
	baseDotfileDir              = ".std"
	baseTestingDir              = ".testing"
	defaultListDatabasePath     = "db/listDB.db"
	defaultPasswordDatabasePath = "db/passwordDB.db"
	defaultTempWorkSpace        = "tmp/"
	defaultLogFile              = "logging/std.log"
	defaultBufferMDFile         = "logging/stdin.buffer.md"
	dotFileConfig               = "config.yml"
)

// STDConf embodies the set configurations for the application
var STDConf = &config{}

var root string

func init() {
	if err := SetConfigWithUserRoot(); err != nil {
		log.Fatal("couldn't set the necessary directory names", err)
	}

	if err := makePaths(); err != nil {
		log.Fatal("couldn't make the neccessary paths", err)
	}

	if err := loadDotfileConfig(); err != nil {
		log.Fatal("Cannot load yaml config:", err)
	}
}

// SetConfigWithUserRoot will set the configs again.
// Run this when running tests
func SetConfigWithUserRoot() error {
	usr, err := userOS.Current()
	if err != nil {
		return err
	}

	root = filepath.Join(usr.HomeDir, baseDotfileDir)
	testingEnvironment, _ := strconv.ParseBool(os.Getenv(testingConfigEnv))
	if testingEnvironment {
		root = filepath.Join(usr.HomeDir, baseDotfileDir, baseTestingDir)
	}

	STDConf.dotFileConfig = filepath.Join(root, dotFileConfig)
	STDConf.BufferMDFile = filepath.Join(root, defaultBufferMDFile)
	STDConf.ListDatabasePath = filepath.Join(root, defaultListDatabasePath)
	STDConf.PasswordDatabasePath = filepath.Join(root, defaultPasswordDatabasePath)
	STDConf.LogFile = filepath.Join(root, defaultLogFile)
	STDConf.TempWorkSpace = filepath.Join(root, defaultTempWorkSpace)
	return nil
}

func makePaths() error {
	if err := os.MkdirAll(filepath.Dir(STDConf.dotFileConfig), 0744); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(STDConf.BufferMDFile), 0744); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(STDConf.ListDatabasePath), 0744); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(STDConf.PasswordDatabasePath), 0744); err != nil {
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
