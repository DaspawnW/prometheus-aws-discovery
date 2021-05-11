package awsdiscovery

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/test"
	"github.com/mitchellh/mapstructure"
	"github.com/nsf/jsondiff"
)

func TestDiscoveryHasCorrectEndpointsV1(t *testing.T) {
	opts := jsondiff.DefaultConsoleOptions()

	scrapeTargets := make([]discovery.Instance, 0)

	for _, instanceData := range test.Instances {
		var instance discovery.Instance
		mapstructure.Decode(instanceData, &instance)
		scrapeTargets = append(scrapeTargets, instance)
	}
	testJSONContentBytes, _ := discovery.TargetConfigBytes(scrapeTargets)

	ec2Client := &MockEC2Client{
		instances: EC2InstanceList(),
		err:       nil,
	}

	d := &DiscoveryClientAWS{
		TagPrefix: true,
		Tag:       "prom/scrape",
	}

	d.SetEC2Client(ec2Client)
	scrapeTargets, err := d.GetInstances()
	if err != nil {
		t.Errorf("Failed to discover instances %v", err)
	}

	scrapeTargetBytes, err := discovery.TargetConfigBytes(scrapeTargets)
	if err != nil {
		t.Errorf("TargetConfigBytes %v", err)
	}
	res, diff := jsondiff.Compare(scrapeTargetBytes, testJSONContentBytes, &opts)
	if res != 0 {
		t.Errorf("JSONDiff V2 Tag \n%v \n%v", res, diff)
	}

}
func TestDiscoveryHasCorrectEndpointsV2(t *testing.T) {
	opts := jsondiff.DefaultConsoleOptions()

	ec2Client := &MockEC2Client{
		instances: EC2InstanceList(),
		err:       nil,
	}

	d := &DiscoveryClientAWS{
		TagPrefix: false,
		Tag:       "prom",
	}
	d.SetEC2Client(ec2Client)

	scrapeTargets, err := d.parseScrapeConfigs(EC2InstanceList())
	if err != nil {
		t.Errorf("parseScrapeConfigsError %v", err)
	}
	if len(scrapeTargets) != 3 {
		t.Errorf("Expected three Instances Targets")
	}
	if len(scrapeTargets[0].Metrics)+len(scrapeTargets[1].Metrics) != 0 {
		t.Errorf("Expected first and Second Instance to not have Targets")
	}
	if len(scrapeTargets[2].Metrics) != 2 {
		t.Errorf("Expected third to have two Targets")
	}

	scrapeTargetBytes, err := discovery.TargetConfigBytes(scrapeTargets)
	if err != nil {

		t.Errorf("TargetConfigBytes %v", err)
	}
	testJSONContentBytes, _ := discovery.TargetConfigBytes(scrapeTargets)

	res, diff := jsondiff.Compare(scrapeTargetBytes, testJSONContentBytes, &opts)
	if res != 0 {
		t.Log("Expected JSON Strings to match")
		t.Errorf("JSONDiff V2 Tag \n%v \n%v", res, diff)
	}

}

type MockEC2Client struct {
	ec2iface.EC2API
	instances []*ec2.Instance
	err       error
}

func (c *MockEC2Client) DescribeInstances(in *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	var reservations []*ec2.Reservation
	var reservation ec2.Reservation

	reservation.Instances = c.instances
	reservations = append(reservations, &reservation)

	return &ec2.DescribeInstancesOutput{
		Reservations: reservations,
	}, nil
}
func EC2InstanceList() []*ec2.Instance {
	var instances []*ec2.Instance

	// instance 1
	instances = append(instances, &ec2.Instance{
		InstanceType:     aws.String("t2.medium"),
		PrivateIpAddress: aws.String("127.0.0.1"),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("prom/scrape:9100/metrics"),
				Value: aws.String("node_exporter"),
			},

			{
				Key:   aws.String("prom/scrape:8080:https/metrics"),
				Value: aws.String("blackbox_exporter"),
			},
			{
				Key:   aws.String("Name"),
				Value: aws.String("Testinstance1"),
			},
			{
				Key:   aws.String("billingnumber"),
				Value: aws.String("1111"),
			},
		},
	})

	// instance 2
	instances = append(instances, &ec2.Instance{
		InstanceType:     aws.String("t2.small"),
		PrivateIpAddress: aws.String("127.0.0.2"),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("prom/scrape:9100/metrics"),
				Value: aws.String("node_exporter"),
			},
			{
				Key:   aws.String("Name"),
				Value: aws.String("Testinstance2"),
			},
			{
				Key:   aws.String("billingnumber"),
				Value: aws.String("2222"),
			},
			{
				Key:   aws.String("prom/scrape:8888"),
				Value: aws.String("test_exporter"),
			},
		},
	})
	//instance v2 tags
	//
	// instance 3
	instances = append(instances, &ec2.Instance{
		InstanceType:     aws.String("t2.small"),
		PrivateIpAddress: aws.String("127.0.0.2"),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("prom"),
				Value: aws.String(test.V2TagValue),
			},
			{
				Key:   aws.String("Name"),
				Value: aws.String("Testinstance3"),
			},
			{
				Key:   aws.String("billingnumber"),
				Value: aws.String("2222"),
			},
		},
	})

	return instances
}
