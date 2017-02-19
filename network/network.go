package network

import (
  "fmt"
  "net"
  "net/http"
  "io/ioutil"
  "regexp"
  "strings"
)

var lookups = []string{
  "http://icanhazip.com",
  "http://bot.whatismyipaddress.com",
  "http://ifconfig.me",
}

/*
  This setup may be a bit verbose for current puproses - I guess I tried to
  implement this in such a way that IPv6 integration could be easily tied in?
  Idk, maybe just a simple string representing the IP would be fine.
*/
type IPAddressString interface {
  NetworkAddress() string
}

type IPFourAddressString struct {
  addr string
}
func (istring *IPFourAddressString) NetworkAddress() (ipAddr string) {
  return istring.addr
}

type IPClassError struct {
  msg string
}
func (err *IPClassError) Error() (msg string) {
  return fmt.Sprintf("IPClassError: %s", err.msg)
}

type InterfaceError struct {
  msg string
}
func (err *InterfaceError) Error() (msg string) {
  return fmt.Sprintf("InterfaceError: %s", err.msg)
}

type NetworkError struct {
  msg string
}
func (err *NetworkError) Error() (msg string) {
  return fmt.Sprintf("NetworkError: %s", msg)
}

/*
  This function gets the network address of a local interface.
  This, like it's sibling function, is used depending on a configuration value -
  and this function specifically is not the default.
*/
func GetLocalNetwork(ifaceName string) (ipAddr IPAddressString, err error) {

  iface, err := net.InterfaceByName(ifaceName)

  if err != nil {
    return nil, err
  }
  addrs, err := iface.Addrs()
  if err != nil {
    return nil, err
  }
  if len(addrs) > 2 {
    return nil, &InterfaceError{msg: fmt.Sprintf("More than two addresses (%d) found on Interface %s.", len(addrs), ifaceName)}
  }

  for _, addr := range addrs {
    switch v := addr.(type) {
    case *net.IPNet:
      if v.IP.To4() != nil {
        ipstr := fmt.Sprintf("%s", v.IP) // string(v.IP) was not working???
        b, err := validateAddress(ipstr); if (b) {
          return &IPFourAddressString{addr: ipstr}, nil
        } else {
          return nil, err
        }
      }
    case *net.IPAddr:
      if v.IP.To4() != nil {
        ipstr := fmt.Sprintf("%s", v.IP)
        b, err := validateAddress(ipstr); if (b) {
          return &IPFourAddressString{addr: ipstr}, nil
        } else {
          return nil, err
        }
      }
    }
  }
  return nil, &InterfaceError{msg: fmt.Sprintf("No valid IPv4 address found - IPv6 is not currently supported if that is in use")}
}

/*
  This function uses the 'lookups' array to discover the network's external IP address.
  These sites return the IP public address of the client request - so we're
  just making http calls here, and the first one that completes successfuly we'll
  use as our IP address.
*/
func GetNatNetwork() (ipAddr IPAddressString, err error) {
  for _, site := range lookups {
    resp, err := http.Get(site)
    if err == nil {
      defer resp.Body.Close()
      body, err := ioutil.ReadAll(resp.Body)
      if err == nil {
        b, err := validateAddress(strings.TrimSpace(string(body))); if (b) {
          return &IPFourAddressString{addr: strings.TrimSpace(string(body))}, nil
        } else {
          return nil, err
        }
      } else {
        return nil, err
      }
    } else {
      return nil, err
    }
  }
  return nil, &NetworkError{msg: "Could not get external IP"}
}

/*
  This is a simple function to do a safety check against the IP that we've
  discovered - if it looks like a privately-classed network address, we don't
  want to do anything with it.
*/
func validateAddress(ipAddr string) (valid bool, err error) {
  var privateAddress = regexp.MustCompile("^(192\\.168)|(172\\.(?:(?:1[6-9])|(?:2[0-9])|(?:3[0-1])))|(10)\\..*")
  // https://golang.org/pkg/regexp/#Regexp.MatchString
  v := privateAddress.MatchString(ipAddr)

  if (v) {
    return false, &IPClassError{msg: fmt.Sprintf("Network Address '%s' looks like a privately-classed network", ipAddr)}
  }
  return true, nil
}
