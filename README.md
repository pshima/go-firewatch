# go-firewatch
A utility to watch sources for signs of fires.  Built for use with Consul (https://consul.io)

Inputs:
- AWS cloudwatch alarms (via prefix name only)

Outputs:
- Exit codes with output (support for consul, nagios, etc)

# Why?

There are a lot of utilities out there that already exist to check CloudWatch metrics.  This tools goals are slightly different.  Majority of the tools or integrations out there directly check CloudWatch metrics and define thresholds in the monitoring.  This requires a fairly tight coupling between changes in infrastructure or applications and monitoring thresholds.

Instead of pushing the logic into the monitoring script, I prefer to keep the CloudWatch alarm configuration inside Terraform.  When my infrastructure changes so does my alarming and the monitoring code does not need to be updated.  This uses the concept of a standardized prefix to allow for a search and check of alarms in a failed state.

If your alarming is built around Consul then this helps with consistency in alarming/notification via Consul, rather than split-brain across Consul and CloudWatch.

The goal here at the moment is only to provide notification of the state of an external alarm to Consul, not to have logic around alarming.

# Usage

```
./go-firewatch -prefix myalarmprefix
```

As a consul alert:
```
{
  "check": {
    "id": "check-cloudwatch-alarms-www",
    "name": "WWW Cloudwatch Alarms - see runbook: https://my.runbooks/www",
    "script": "/path/to/go-firewatch -prefix prod-www-alarms",
    "interval": "60s",
    "timeout": "10s"
  }
}
```

# Examples
```
% ./go-firewatch -prefix test
2016/06/23 17:45:33 Error: testing123: CPUUtilization, Threshold Crossed: 1 datapoint (1.554) was greater than the threshold (0.0).
% echo $?
255
```

```
% ./go-firewatch -prefix foo
% echo $?
0
```
