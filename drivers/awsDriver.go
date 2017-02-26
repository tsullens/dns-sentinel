package drivers

import (
  "errors"
  "fmt"
  "log"
  "github.com/aws/aws-sdk-go/service/route53"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/spf13/viper"
)

// Info
// HostedZone https://docs.aws.amazon.com/sdk-for-go/api/service/route53/#HostedZone

type AwsDriver struct {
  ZoneName        string
  RecordName      string
  ProviderConfig  *viper.Viper
  AppLogger       *log.Logger
}

const RRType = "A"

func (driver *AwsDriver) Run(ipAddr string) {

  sess, err := session.NewSession()
  if err != nil {
    driver.AppLogger.Printf("Failed to created session: %s", err)
    return
  }
  svc := route53.New(sess)
  awsHostedZone, err := getHostedZoneId(svc, driver.ZoneName)
  if err != nil {
    driver.AppLogger.Println(err.Error())
    return
  }
  rRecordSet, err := getRecordSet(svc, awsHostedZone, driver.RecordName)
  if err != nil {
    driver.AppLogger.Println(err)
    return
  }
  if *rRecordSet.Type == RRType {
    if *rRecordSet.ResourceRecords[0].Value != ipAddr {
      change := formatChange(driver.RecordName, ipAddr)
      driver.AppLogger.Printf("Submitting change for %s : %s -> %s\n", driver.RecordName, *rRecordSet.ResourceRecords[0].Value, ipAddr)
      resp, err := submitChanges(svc, []*route53.Change{change}, awsHostedZone)
      if err != nil {
        driver.AppLogger.Println(err.Error())
        return
      } else {
        driver.AppLogger.Println("ChangeSet submitted successfully: %s", resp)
      }
    }
  }
  driver.AppLogger.Println("No change needed.")
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
  return nil, errors.New(fmt.Sprintf("Hosted Zone '%s' not found", zoneName))
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
  return nil, errors.New(fmt.Sprintf("Failed to find record %s in zone %s", recordName, *awsHostedZone.Name))
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
  resp, err = svc.ChangeResourceRecordSets(reqParams)
  if err != nil {
    return nil, err
  }
  return resp, nil
}
