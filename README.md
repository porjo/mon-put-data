
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

* Run command:

```
mon-put-data -metric-name "processCount" -namespace "CustomMetrics" -resolution 60 -unit "Count" -value 12
```

- `-resolution` can be either `1` for high resolution or `60` for standard resolution
- `-unit` can be any of:
```
Seconds | Microseconds | Milliseconds | Bytes | Kilobytes | Megabytes | Gigabytes | Terabytes | Bits | Kilobits | Megabits | Gigabits | Terabits | Percent | Count | Bytes/Second | Kilobytes/Second | Megabytes/Second | Gigabytes/Second | Terabytes/Second | Bits/Second | Kilobits/Second | Megabits/Second | Gigabits/Second | Terabits/Second | Count/Second | None
```
(Ref: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_MetricDatum.html)

## Credit

Code based on https://github.com/mlabouardy/mon-put-instance-data
