package main

import (
  "fmt"
  "os"
  "log"
  "flag"
  "dns-sentinel/drivers"
  "github.com/spf13/viper"
)

type RunConfig struct {
  RunInterval int
  NetFunc   func(cfg *viper.Viper) (ipAddr IPAddressString, err error)
  Provider  drivers.Driver
}

var (
  app_version     string = "0.0.1"
  Runtime_Logger  *log.Logger
  Runtime_Config  *RunConfig
  runtime_viper   *viper.Viper
)

func Init() {

  setupConfig()
  validateAndLoadConfig()
}

func setupConfig() {
  runtime_viper = viper.New()

  // Add flag to specify config file

  runtime_viper.SetConfigType("toml")
  runtime_viper.SetConfigName("test")
  runtime_viper.AddConfigPath(".")
  runtime_viper.AddConfigPath("/etc/dns-sentinel")

  runtime_viper.SetDefault("appconfig", map[string]interface{}{"log": "/var/log/dns-sentinel.log", "loglevel": "info", "run_interval": 1800})
  runtime_viper.SetDefault("network.type", "nat")
}

func validateAndLoadConfig() {

  err :=  runtime_viper.ReadInConfig()
  if err != nil {
    log.Fatalf(fmt.Sprintf("Fatal error reading config file: %s \n", err))
  }

  lfile, err := os.Create(runtime_viper.GetString("appconfig.logfile"))
  if err != nil {
    log.Fatalf("%s\n", err.Error())
  }
  Runtime_Logger = log.New(lfile, "", log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

  // Start bulding our running config struct
  Runtime_Config = &RunConfig{
    RunInterval:  runtime_viper.GetInt("appconfig.run_interval"),
  }

  switch runtime_viper.GetString("zone.provider") {
  case "aws":
    if runtime_viper.IsSet("zone.name") && runtime_viper.IsSet("zone.record") {
      Runtime_Config.Provider = &drivers.AwsDriver{
        ZoneName:        fmt.Sprintf("%s.", runtime_viper.GetString("zone.name")),
        RecordName:      fmt.Sprintf("%s.", runtime_viper.GetString("zone.record")),
        ProviderConfig:  runtime_viper.Sub("provider.aws"),
        AppLogger:       Runtime_Logger,
      }
    } else {
      Runtime_Logger.Fatalf("Missing required parameters for 'zone.provider': '%s'\n", runtime_viper.GetString("zone.provider"))
    }
  default:
    Runtime_Logger.Fatalf("%s is Not a valid 'provider'.\n", runtime_viper.GetString("zone.provider"))
  }

  switch runtime_viper.GetString("network.type") {
  case "local":
    if viper.IsSet("network.interface") {
      Runtime_Config.NetFunc = GetLocalNetwork
    } else {
      Runtime_Logger.Fatalf("Network Config Error: 'interface' value is required for a 'type' of 'local'.\n")
    }
  case "nat":
    Runtime_Config.NetFunc = GetNatNetwork
  case "default":
    Runtime_Logger.Fatalf("Network Config Error: '%s' is not a valid value for 'type'.\n")
  }
  Runtime_Logger.Printf("Succesfully loaded config: %v\n", runtime_viper)
}
