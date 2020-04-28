package endpoints

import (
	"testing"
)

func TestToJsonStringHasCorrectContent(t *testing.T) {

	instanceList := InstanceList()

	returnedJSONContentBytes, _ := ToJsonString(DiscoveredInstances{Instances: instanceList})
	expectedJSONContent := "[{\"targets\":[\"127.0.0.1:9100\"],\"labels\":{\"__address__\":\"127.0.0.1:9100\",\"__metrics_path__\":\"/metrics\",\"billingnumber\":\"1111\",\"instancename\":\"Testinstance1\",\"name\":\"node_exporter\"}},{\"targets\":[\"127.0.0.1:8080\"],\"labels\":{\"__address__\":\"127.0.0.1:8080\",\"__metrics_path__\":\"/metrics\",\"billingnumber\":\"1111\",\"instancename\":\"Testinstance1\",\"name\":\"blackbox_exporter\"}},{\"targets\":[\"127.0.0.2:9100\"],\"labels\":{\"__address__\":\"127.0.0.2:9100\",\"__metrics_path__\":\"/metrics\",\"billingnumber\":\"2222\",\"instancename\":\"Testinstance2\",\"name\":\"node_exporter\"}}]"

	if expectedJSONContent != string(returnedJSONContentBytes) {
		t.Errorf("Expected json string with content %s, but got %s", expectedJSONContent, string(returnedJSONContentBytes))
	}

}

func InstanceList() []Instance {
	expectedInstanceList := []Instance{}

	// Instance 1
	tagsI1 := make(map[string]string)
	tagsI1["Name"] = "Testinstance1"
	tagsI1["billingnumber"] = "1111"
	metricsI1 := []InstanceMetrics{}
	m1I1 := InstanceMetrics{
		Name: "node_exporter",
		Path: "/metrics",
		Port: 9100,
	}
	m2I1 := InstanceMetrics{
		Name: "blackbox_exporter",
		Path: "/metrics",
		Port: 8080,
	}
	metricsI1 = append(metricsI1, m1I1, m2I1)
	i1 := Instance{
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
	metricsI2 := []InstanceMetrics{}
	m1I2 := InstanceMetrics{
		Name: "node_exporter",
		Path: "/metrics",
		Port: 9100,
	}
	metricsI2 = append(metricsI2, m1I2)
	i2 := Instance{
		InstanceType: "t2.small",
		PrivateIP:    "127.0.0.2",
		Tags:         tagsI2,
		Metrics:      metricsI2,
	}
	expectedInstanceList = append(expectedInstanceList, i2)
	return expectedInstanceList
}
