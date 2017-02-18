package drivers

import (
  "errors"
  "fmt"
  "os"
  "github.com/aws/aws-sdk-go/service/route53"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
)

// Info
// HostedZone https://docs.aws.amazon.com/sdk-for-go/api/service/route53/#HostedZone

type AwsDriver struct {
  IpAddr string
  ZoneName string
  RecordName string
}

const RRType = "A"

func (driver *AwsDriver) Run() {

  sess, err := session.NewSession()
  if err != nil {
    fmt.Println("Failed to created session, ", err)
    os.Exit(1)
  }
  svc := route53.New(sess)
  awsHostedZone, err := getHostedZoneId(svc, driver.ZoneName)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }
  rRecordSet, err := getRecordSet(svc, awsHostedZone, driver.RecordName)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  if *rRecordSet.Type == RRType {
    if *rRecordSet.ResourceRecords[0].Value != driver.IpAddr {
      change := formatChange(driver.RecordName, driver.IpAddr)
      fmt.Printf("Submitting change for %s : %s -> %s\n", driver.RecordName, *rRecordSet.ResourceRecords[0].Value, driver.IpAddr)
      resp, err := submitChanges(svc, []*route53.Change{change}, awsHostedZone)
      if err != nil {
        fmt.Println(err.Error())
      } else {
        fmt.Println(resp)
      }
    }
  }
  fmt.Println("No change needed.")
}

func getHostedZoneId(svc *route53.Route53, zoneName string) (awsHostedZone *route53.HostedZone, err error ){
  reqParams := &route53.ListHostedZonesByNameInput{}
  resp, err := svc.ListHostedZonesByName(reqParams)
  if err != nil {
    return nil, err
  }
  for _, hostedZone := range resp.HostedZones {
    if *hostedZone.Name == zoneName {
      return hostedZone, nil
    }
  }
  return nil, errors.New("Not Found")
}

func getRecordSet(svc *route53.Route53, awsHostedZone *route53.HostedZone, recordName string) (set *route53.ResourceRecordSet, err error) {
  reqParams := &route53.ListResourceRecordSetsInput{
    HostedZoneId: awsHostedZone.Id,
  }
  resp, err := svc.ListResourceRecordSets(reqParams)
  if err != nil {
    return nil, err
  }
  for _, recordSet := range resp.ResourceRecordSets {
    if *recordSet.Name == recordName {
      return recordSet, nil
    }
  }
  return nil, errors.New(fmt.Sprintf("Failed to find record %s in zone %s", recordName, awsHostedZone.Name))
}


// Helper func to create a route53.Change to be submitted to the submitChanges func
func formatChange(recordName, rData string) (change *route53.Change) {

  return &route53.Change{
    Action: aws.String("UPSERT"),
    ResourceRecordSet: &route53.ResourceRecordSet{
      Name: aws.String(recordName),
      Type: aws.String(RRType),
      ResourceRecords: []*route53.ResourceRecord{
        {
          Value: aws.String(rData),
        },
      },
      TTL: aws.Int64(300),
    },
  }
}

func submitChanges(svc *route53.Route53, changes []*route53.Change, awsHostedZone *route53.HostedZone) (resp *route53.ChangeResourceRecordSetsOutput, err error) {
  reqParams := &route53.ChangeResourceRecordSetsInput{
    ChangeBatch: &route53.ChangeBatch{
      Changes: changes,
    },
    HostedZoneId: awsHostedZone.Id,
  }
  fmt.Println(reqParams)
  resp, err = svc.ChangeResourceRecordSets(reqParams)
  if err != nil {
    return nil, err
  }
  return resp, nil
}
