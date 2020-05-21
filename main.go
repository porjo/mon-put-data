package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	//	"github.com/aws/aws-sdk-go/aws/external"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type dimensions []*cloudwatch.Dimension

func (i *dimensions) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *dimensions) Set(value string) error {
	bits := strings.SplitN(value, "=", 2)
	if len(bits) == 2 {
		d := &cloudwatch.Dimension{}
		d.Name = &bits[0]
		d.Value = &bits[1]
		*i = append(*i, d)
	}
	return nil
}

func main() {

	var dims dimensions

	region := flag.String("region", "ap-southeast-2", "")
	namespace := flag.String("namespace", "CustomMetrics", "")
	metricName := flag.String("metric-name", "", "the name of the metric")
	unit := flag.String("unit", "Count", "Cloudwatch metric unit")
	value := flag.Float64("value", 0, "")
	resolution := flag.Int64("resolution", 60, "storage resolution for metric in seconds (1 or 60)")
	flag.Var(&dims, "dimension", "name=value pair that uniquely identifies metric. Multiple allowed. Defaults to instanceID")
	flag.Parse()

	if *metricName == "" {
		fmt.Printf("Please supply metric-name\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *resolution != 1 && *resolution != 60 {
		fmt.Printf("resolution must be either 1 or 60\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg := aws.NewConfig()
	cfg.Region = region
	cloudWatch := CloudWatchService{
		Config: cfg,
	}

	if dims == nil {
		id, err := GetInstanceID()
		if err != nil {
			log.Fatal(err)
		}
		dimensionKey := "InstanceId"
		dims = []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  &dimensionKey,
				Value: &id,
			},
		}
	}
	metricDatum := constructMetricDatum(*metricName, *value, *resolution, *unit, dims)
	cloudWatch.Publish(metricDatum, *namespace)
}

// CloudWatchService entity
type CloudWatchService struct {
	Config *aws.Config
}

// Publish save metrics to cloudwatch using AWS CloudWatch API
func (c CloudWatchService) Publish(metricData []*cloudwatch.MetricDatum, namespace string) {
	mySession := session.Must(session.NewSession())
	svc := cloudwatch.New(mySession, c.Config)
	req, _ := svc.PutMetricDataRequest(&cloudwatch.PutMetricDataInput{
		MetricData: metricData,
		Namespace:  &namespace,
	})
	err := req.Send()
	if err != nil {
		log.Fatal(err)
	}
}

// constructMetricDatum construct cloudwatch data object
func constructMetricDatum(metricName string, value float64, resolution int64, unit string, dimensions []*cloudwatch.Dimension) []*cloudwatch.MetricDatum {
	datum := &cloudwatch.MetricDatum{
		MetricName:        &metricName,
		Dimensions:        dimensions,
		Unit:              &unit,
		Value:             &value,
		StorageResolution: &resolution,
	}
	return []*cloudwatch.MetricDatum{datum}
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
