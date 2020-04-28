package libtesting

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/endpoints"
)

func InstanceList() []endpoints.Instance {
	expectedInstanceList := []endpoints.Instance{}

	// Instance 1
	tagsI1 := make(map[string]string)
	tagsI1["Name"] = "Testinstance1"
	tagsI1["billingnumber"] = "1111"
	metricsI1 := []endpoints.InstanceMetrics{}
	m1I1 := endpoints.InstanceMetrics{
		Name: "node_exporter",
		Path: "/metrics",
		Port: 9100,
	}
	m2I1 := endpoints.InstanceMetrics{
		Name: "blackbox_exporter",
		Path: "/metrics",
		Port: 8080,
	}
	metricsI1 = append(metricsI1, m1I1, m2I1)
	i1 := endpoints.Instance{
		InstanceType: "t2.medium",
		PrivateIP:    "127.0.0.1",
		Tags:         tagsI1,
		Metrics:      metricsI1,
	}
	expectedInstanceList = append(expectedInstanceList, i1)

	// Instance 2
	tagsI2 := make(map[string]string)
	tagsI2["Name"] = "Testinstance2"
	tagsI2["billingnumber"] = "2222"
	metricsI2 := []endpoints.InstanceMetrics{}
	m1I2 := endpoints.InstanceMetrics{
		Name: "node_exporter",
		Path: "/metrics",
		Port: 9100,
	}
	metricsI2 = append(metricsI2, m1I2)
	i2 := endpoints.Instance{
		InstanceType: "t2.small",
		PrivateIP:    "127.0.0.2",
		Tags:         tagsI2,
		Metrics:      metricsI2,
	}
	expectedInstanceList = append(expectedInstanceList, i2)
	return expectedInstanceList
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
		Key:   aws.String("prom/scrape:8080/metrics"),
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
