
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
mon-put-data --metric-name "processCount" --namespace "CustomMetrics" --unit "Count" --value 12
```

## Credit

Code based on https://github.com/mlabouardy/mon-put-instance-data
