package firewatch

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type Job struct {
	Prefix    string
	RawOutput []*cloudwatch.MetricAlarm
}

func (j *Job) Check() error {
	out, err := queryCloudWatchAlarm(j.Prefix)
	if err != nil {
		return err
	}
	j.RawOutput = out
	if len(out) == 0 {
		return nil
	} else {
		return fmt.Errorf("%s: %s, %s",
			*j.RawOutput[0].AlarmName,
			*j.RawOutput[0].MetricName,
			*j.RawOutput[0].StateReason)
	}

	return nil
}

func queryCloudWatchAlarm(p string) ([]*cloudwatch.MetricAlarm, error) {
	region, err := getRegion()
	if err != nil {
		return nil, err
	}

	svc := cloudwatch.New(session.New(&aws.Config{Region: aws.String(region)}))

	params := &cloudwatch.DescribeAlarmsInput{
		AlarmNamePrefix: aws.String(p),
		StateValue:      aws.String("ALARM"),
	}

	resp, err := svc.DescribeAlarms(params)
	if err != nil {
		return nil, err
	}

	return resp.MetricAlarms, nil
}
