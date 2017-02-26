package main

import (
  "time"
)

func main() {

  setupConfig()
  validateAndLoadConfig()

  for {
    Runtime_Logger.Printf("Starting Run")
    ipAddr, err := Runtime_Config.NetFunc(runtime_viper.Sub("network"))
    if err != nil {
      Runtime_Logger.Printf("Error resolving IP Address: %s", err.Error())
    } else {
      Runtime_Logger.Printf("Found IP Address: %s", ipAddr.NetworkAddress())
      Runtime_Config.Provider.Run(ipAddr.NetworkAddress())
    }
    Runtime_Logger.Printf("Sleeping for %d seconds\n", Runtime_Config.RunInterval)
    time.Sleep(time.Duration(Runtime_Config.RunInterval) * time.Second)
  }
}
