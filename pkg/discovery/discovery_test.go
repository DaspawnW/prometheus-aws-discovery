package discovery

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/libtesting"
)

func TestDiscoveryHasCorrectEndpoints(t *testing.T) {
	ec2Client := &MockEC2Client{
		instances: libtesting.EC2InstanceList(),
		err:       nil,
	}
	d := NewDiscovery(ec2Client, "prom/scrape")
	returnedInstanceList, err := d.Discover()
	if err != nil {
		t.Error("Failed to discover instances", err)
	}

	expectedInstanceList := libtesting.InstanceList()

	if !reflect.DeepEqual(expectedInstanceList, returnedInstanceList.Instances) {
		t.Errorf("Expected instance list %v to equal returned instance list %v", expectedInstanceList, returnedInstanceList)
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
