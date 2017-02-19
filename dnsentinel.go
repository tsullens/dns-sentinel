package main

import (
  "fmt"
  "os"
  "dnsentinel/network"
  "dnsentinel/drivers"
)

const (
  domainName = "homelab.tcsullens.com." // Zone Name
  discovery = "nat"
  ifaceName = "en0"
  recordName = "vpn"
  validRecord = "[a-zA-Z]+(\\.[a-zA-Z]+)+\\.?"
)

func main() {

  ipAddr, err := network.GetLocalNetwork(ifaceName)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
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
