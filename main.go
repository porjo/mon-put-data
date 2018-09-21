package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

func main() {

	region := flag.String("region", "ap-southeast-2", "")
	namespace := flag.String("namespace", "CustomMetrics", "")
	metricName := flag.String("metric-name", "", "")
	unit := flag.String("unit", "Count", "Cloudwatch metric unit")
	value := flag.Float64("value", 0, "")
	flag.Parse()

	if *metricName == "" {
		fmt.Printf("Please supply metric-name\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("Unable to load SDK config")
	}
	cfg.Region = *region
	cloudWatch := CloudWatchService{
		Config: cfg,
	}

	id, err := GetInstanceID()
	if err != nil {
		log.Fatal(err)
	}

	dimensionKey := "InstanceId"
	dimensions := []cloudwatch.Dimension{
		cloudwatch.Dimension{
			Name:  &dimensionKey,
			Value: &id,
		},
	}
	metricDatum := constructMetricDatum(*metricName, *value, cloudwatch.StandardUnit(*unit), dimensions)
	cloudWatch.Publish(metricDatum, *namespace)
}

// CloudWatchService entity
type CloudWatchService struct {
	Config aws.Config
}

// Publish save metrics to cloudwatch using AWS CloudWatch API
func (c CloudWatchService) Publish(metricData []cloudwatch.MetricDatum, namespace string) {
	svc := cloudwatch.New(c.Config)
	req := svc.PutMetricDataRequest(&cloudwatch.PutMetricDataInput{
		MetricData: metricData,
		Namespace:  &namespace,
	})
	_, err := req.Send()
	if err != nil {
		log.Fatal(err)
	}
}

// constructMetricDatum construct cloudwatch data object
func constructMetricDatum(metricName string, value float64, unit cloudwatch.StandardUnit, dimensions []cloudwatch.Dimension) []cloudwatch.MetricDatum {
	return []cloudwatch.MetricDatum{
		cloudwatch.MetricDatum{
			MetricName: &metricName,
			Dimensions: dimensions,
			Unit:       unit,
			Value:      &value,
		},
	}
}

// GetInstanceID return EC2 instance id
func GetInstanceID() (string, error) {
	value := os.Getenv("AWS_INSTANCE_ID")
	if len(value) > 0 {
		return value, nil
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://169.254.169.254/latest/meta-data/instance-id", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
