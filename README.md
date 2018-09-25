
# mon-put-data

Update a single Cloudwatch metric using a single binary (no dependencies).

## How to use

* Setup an IAM Policy:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "1",
            "Effect": "Allow",
            "Action": "cloudwatch:PutMetricData",
            "Resource": "*"
        }
    ]
}
```

* Create an IAM role using the above policy, and attach it to your ec2 instance

* Run command

```
$ mon-put-data

  -dimension value
    	name=value pair that uniquely identifies metric. Multiple allowed. Defaults to instanceID
  -metric-name string
    	the name of the metric
  -namespace string
    	 (default "CustomMetrics")
  -region string
    	 (default "ap-southeast-2")
  -resolution int
    	storage resolution for metric in seconds (1 or 60) (default 60)
  -unit string
    	Cloudwatch metric unit (default "Count")
  -value float
```

- `-unit` can be any of:
```
Seconds | Microseconds | Milliseconds | Bytes | Kilobytes | Megabytes | Gigabytes | Terabytes | Bits | Kilobits | Megabits | Gigabits | Terabits | Percent | Count | Bytes/Second | Kilobytes/Second | Megabytes/Second | Gigabytes/Second | Terabytes/Second | Bits/Second | Kilobits/Second | Megabits/Second | Gigabits/Second | Terabits/Second | Count/Second | None
```
(Ref: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_MetricDatum.html)

## Example Usage

Want to monitor the average process count for a group of Apache servers. Each server would run this command every 60 seconds:

```
mon-put-data -namespace "CustomMetrics" -dimension "Apache Stats=Processes" -metric-name "processCount" -resolution 60 -unit "Count" -value $COUNT
```

In Cloudwatch an alarm can be created using this metric and specifying 'Average' and '1 Minute' time threshold parameters.

## Credit

Code based on https://github.com/mlabouardy/mon-put-instance-data
