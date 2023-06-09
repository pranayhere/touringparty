package config

import (
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// Default options for configuration loading
const (
	DefaultConfigType     = "yaml"
	DefaultConfigDir      = "./config"
	DefaultConfigFileName = "default"
	WorkDirEnv            = "WORKDIR"
)

// Options is config options
type Options struct {
	configType            string
	configPath            string
	defaultConfigFileName string
}

// Config is a wrapper over a underlying config loader implementation
type Config struct {
	opts  Options
	viper *viper.Viper
}

// NewDefaultOptions returns default options.
// DISCLAIMER: This function is a bit hacky
// This function expects an env $WORKDIR to
// be set and reads configs from $WORKDIR/configs.
// If $WORKDIR is not set. It uses the absolute path wrt
// the location of this file (config.go) to set configPath
// to 2 levels up in viper (../../configs).
// This function breaks if :
// 1. $WORKDIR is set and configs dir not present in $WORKDIR
// 2. $WORKDIR is not set and ../../configs is not present
// 3. $WORKDIR is not set and runtime absolute path of configs
// is different than build time path as runtime.Caller() evaluates
// only at build time
func NewDefaultOptions() Options {
	var configPath string
	workDir := os.Getenv(WorkDirEnv)
	if workDir != "" {
		configPath = path.Join(workDir, DefaultConfigDir)
	} else {
		_, thisFile, _, _ := runtime.Caller(1)
		configPath = path.Join(path.Dir(thisFile), "../../"+DefaultConfigDir)
	}
	return NewOptions(DefaultConfigType, configPath, DefaultConfigFileName)
}

// NewOptions returns new Options struct.
func NewOptions(configType string, configPath string, defaultConfigFileName string) Options {
	return Options{configType, configPath, defaultConfigFileName}
}

// NewDefaultConfig returns new config struct with default options.
func NewDefaultConfig() *Config {
	return NewConfig(NewDefaultOptions())
}

// NewConfig returns new config struct.
func NewConfig(opts Options) *Config {
	return &Config{opts, viper.New()}
}

// Load reads environment specific configurations and along with the defaults
// unmarshalls into config.
func (c *Config) Load(env string, config interface{}) error {
	log.Println(c.opts.defaultConfigFileName)
	return c.loadByConfigName(env, config)
}

// loadByConfigName reads configuration from file and unmarshalls into config.
func (c *Config) loadByConfigName(configName string, config interface{}) error {
	c.viper.SetEnvPrefix(strings.ToUpper("touringparty"))
	c.viper.SetConfigName(configName)
	c.viper.SetConfigType(c.opts.configType)
	c.viper.AddConfigPath(c.opts.configPath)
	c.viper.AutomaticEnv()

	log.Println(c.opts.configPath, configName)

	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	return c.viper.Unmarshal(config)
}
