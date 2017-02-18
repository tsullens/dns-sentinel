package network

import (
  "fmt"
  "net"
  "net/http"
  "errors"
  "io/ioutil"
  "strings"
)

var lookups = []string{ "http://icanhazip.com", "http://bot.whatismyipaddress.com", "http://ifconfig.me" }

type NetInterface struct {
  IP_four string
  IP_six string
}

func GetLocalNetwork(ifaceName string) (niface *NetInterface, err error) {

  iface, err := net.InterfaceByName(ifaceName)

  if err != nil {
    return nil, err
  }
  addrs, err := iface.Addrs()
  if err != nil {
    return nil, err
  }
  if len(addrs) > 2 {
    return nil, errors.New(fmt.Sprintf("More than two addresses (%d) found on Interface %s.", len(addrs), ifaceName))
  }

  var ni = &NetInterface{}
  for _, addr := range addrs {
    switch v := addr.(type) {
    case *net.IPNet:
      if v.IP.To4() != nil {
        ni.IP_four = string(v.IP)
      } else {
        ni.IP_six = string(v.IP)
      }
    case *net.IPAddr:
      if v.IP.To4() != nil {
        ni.IP_four = string(v.IP)
      } else {
        ni.IP_six = string(v.IP)
      }
    }
  }
  return ni, nil
}

func GetNatNetwork() (niface *NetInterface, err error) {
  for _, site := range lookups {
    resp, err := http.Get(site)
    if err == nil {
      defer resp.Body.Close()
      body, err := ioutil.ReadAll(resp.Body)
      if err == nil {
        return &NetInterface{IP_four: strings.TrimSpace(string(body)), IP_six: ""}, nil
      } else {
        fmt.Println(err.Error())
      }
    } else {
      fmt.Println(err.Error())
    }
  }
  return nil, errors.New("getNatNetwork: Could not get external IP")
}
