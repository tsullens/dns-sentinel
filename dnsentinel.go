package main

import (
  "flag"
  "fmt"
  "os"
  "dnsentinel/network"
  "dnsentinel/drivers"
)

var (
  ifaceName string
  zoneFlag string
  recordFlag string
)

const validRecord = "[a-zA-Z]+(\\.[a-zA-Z]+)+\\.?"

func main() {

  flag.StringVar(&ifaceName, "iface", "eth0", "Name of the interface to watch.")
  flag.StringVar(&zoneFlag, "zone", "", "Hosted Zone. e.g. example.com")
  flag.StringVar(&recordFlag, "record", "", "Record Name")
  flag.Parse()

  niface, err := network.GetNatNetwork()
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }
  //fmt.Println(niface)
  awsDriver := &drivers.AwsDriver{
    IpAddr: niface.IP_four,
    ZoneName: zoneFlag,
    RecordName: recordFlag,
  }
  awsDriver.Run()
}
