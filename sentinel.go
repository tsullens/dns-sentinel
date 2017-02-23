package main

import (
  "fmt"
  "dns-sentinel/network"
  "dns-sentinel/drivers"
  "github.com/spf13/viper"
)

const (
  domainName = "homelab.tcsullens.com." // Zone Name
  discovery = "nat"
  ifaceName = "en0"
  recordName = "vpn"
  validRecord = "[a-zA-Z]+(\\.[a-zA-Z]+)+\\.?"
)

var runtime_viper *viper.Viper
var ipAddr network.IPAddressString
var err error

func main() {

  loadConfig()

  switch runtime_viper.GetString("network.type") {
  case "local":
    //fmt.Printf("%s\n", runtime_viper.GetString("network.interface"))
    if viper.IsSet("network.interface") {
      ipAddr, err = network.GetLocalNetwork(ifaceName)
      if err != nil {
        fmt.Println(err.Error())
      }
    } else {
      fmt.Println("network.interface is required for a network.type of local")
    }
  case "nat":
    //fmt.Printf("NAT\n")
    ipAddr, err = network.GetNatNetwork()
    if err != nil {
      fmt.Println(err.Error())
    }
  case "default":
    fmt.Printf("Not a Valid network type.")
  }

  switch runtime_viper.GetString("provider") {
  case "aws":
    fmt.Printf("%v\n", runtime_viper.GetString("provider"))
  default:
    fmt.Printf("%v is Not a valid provider\n", runtime_viper.GetString("provider"))
  }

  fmt.Println(ipAddr.NetworkAddress())
  awsDriver := &drivers.AwsDriver{
    IpAddr: ipAddr.NetworkAddress(),
    ZoneName: domainName,
    RecordName: fmt.Sprintf("%s.%s", recordName, domainName),
  }
  fmt.Println(awsDriver)
  //awsDriver.Run()
}

func loadConfig() {
  runtime_viper = viper.New()
  // Add flag to specify config file 

  runtime_viper.SetConfigType("toml")
  runtime_viper.SetConfigName("testcfg")
  runtime_viper.AddConfigPath(".")
  runtime_viper.AddConfigPath("/etc/dns-sentinel")

  runtime_viper.SetDefault("appconfig", map[string]interface{}{"log": "/var/log/dns-sentinel.log", "loglevel": "info", "poll_interval": 1800})
  runtime_viper.SetDefault("network.type", "nat")

  err :=  runtime_viper.ReadInConfig()
  if err != nil {
    panic(fmt.Errorf("Fatal error reading config file: %s \n", err))
  }

}
