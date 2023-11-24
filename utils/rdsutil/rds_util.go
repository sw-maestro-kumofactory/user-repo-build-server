package rdsutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSClient struct {
	client *rds.RDS
}

func NewRDSClient() (*RDSClient, error) {
	accessKey := "AKIAWKE7NUBXC47IZYXL"
	secretKey := ""
	region := "ap-northeast-2"

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})

	rdsClient := rds.New(sess)
	return &RDSClient{client: rdsClient}, nil
}

func (rc *RDSClient) GetEndpointByDBName(dbName string) (string, error) {
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbName),
	}

	result, err := rc.client.DescribeDBInstances(input)
	if err != nil {
		return "", err
	}

	if len(result.DBInstances) > 0 {
		endpoint := aws.StringValue(result.DBInstances[0].Endpoint.Address)
		return endpoint, nil
	}

	return "", fmt.Errorf("RDS instance not found")
}

func (rc *RDSClient) GetRdsTypeByDBName(dbName string) (string, error) {
	input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbName),
	}

	result, err := rc.client.DescribeDBInstances(input)
	if err != nil {
		return "", err
	}

	if len(result.DBInstances) > 0 {
		rdsType := aws.StringValue(result.DBInstances[0].Engine)
		return rdsType, nil
	}

	return "", fmt.Errorf("RDS instance not found")
}
