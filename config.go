package main

import (
  "github.com/spf13/viper"
  "fmt"
)

type Config struct {
  AppConfig map[string]string
  Provider map[string]interface{}
  Network map[string]string
  Zone map[string]string
}

func loadConfig() (cfg *Config) {
  viper.SetConfigType("toml")
  viper.SetConfigName("testcfg")
  viper.AddConfigPath(".")
  viper.AddConfigPath("/etc/dns-sentinel")

  viper.SetDefault("appconfig", map[string]string{"log": "/var/log/dns-sentinel.log", "loglevel": "info", "poll_interval": "1800"})
  viper.SetDefault("network.type", "nat")

  err :=  viper.ReadInConfig()
  if err != nil {
    panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
  }

  var cfg Config
  err = viper.Unmarshal(&c)
  if err != nil {
    fmt.Printf("Unable to decode config to struct, %v\n", err)
  }
  return &cfg

}
