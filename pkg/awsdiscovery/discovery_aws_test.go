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

func TestDiscoveryHasCorrectEndpoints(t *testing.T) {

	instances := make([]discovery.Instance, 0)
	for _, instanceData := range test.Instances {
		var instance discovery.Instance
		mapstructure.Decode(instanceData, &instance)
		instances = append(instances, instance)
	}
	testJSONContentBytes, _ := discovery.TargetConfigBytes(instances)

	ec2Client := &MockEC2Client{
		instances: EC2InstanceList(),
		err:       nil,
	}

	d := &DiscoveryClientAWS{
		TagPrefix: "prom/scrape",
	}
	d.SetEC2Client(ec2Client)
	returnedInstanceList, err := d.GetInstances()
	if err != nil {
		t.Error("Failed to discover instances", err)
	}

	returnedJSONContentBytes, err := discovery.TargetConfigBytes(returnedInstanceList)
	opts := jsondiff.DefaultConsoleOptions()
	res, text := jsondiff.Compare(returnedJSONContentBytes, testJSONContentBytes, &opts)
	if res != 0 {
		t.Errorf("JSONDiff \n%v \n%v", res, text)
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
	var tagsInstance1 []*ec2.Tag
	t11 := ec2.Tag{
		Key:   aws.String("prom/scrape:9100/metrics"),
		Value: aws.String("node_exporter"),
	}
	t12 := ec2.Tag{
		Key:   aws.String("prom/scrape:8080:https/metrics"),
		Value: aws.String("blackbox_exporter"),
	}
	nameTag1 := ec2.Tag{
		Key:   aws.String("Name"),
		Value: aws.String("Testinstance1"),
	}
	additionalTag1 := ec2.Tag{
		Key:   aws.String("billingnumber"),
		Value: aws.String("1111"),
	}
	tagsInstance1 = append(tagsInstance1, &t11, &t12, &nameTag1, &additionalTag1)
	i1 := ec2.Instance{
		InstanceType:     aws.String("t2.medium"),
		PrivateIpAddress: aws.String("127.0.0.1"),
		Tags:             tagsInstance1,
	}
	instances = append(instances, &i1)

	// instance 2
	var tagsInstance2 []*ec2.Tag
	t21 := ec2.Tag{
		Key:   aws.String("prom/scrape:9100/metrics"),
		Value: aws.String("node_exporter"),
	}
	nameTag2 := ec2.Tag{
		Key:   aws.String("Name"),
		Value: aws.String("Testinstance2"),
	}
	additionalTag2 := ec2.Tag{
		Key:   aws.String("billingnumber"),
		Value: aws.String("2222"),
	}
	exceptionTag := ec2.Tag{
		Key:   aws.String("prom/scrape:8888"),
		Value: aws.String("test_exporter"),
	}
	tagsInstance2 = append(tagsInstance2, &t21, &nameTag2, &additionalTag2, &exceptionTag)
	i2 := ec2.Instance{
		InstanceType:     aws.String("t2.small"),
		PrivateIpAddress: aws.String("127.0.0.2"),
		Tags:             tagsInstance2,
	}
	instances = append(instances, &i2)

	return instances
}
