package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Instance struct {
	InstanceType string
	PrivateIP    string
	Tags         map[string]string
	Metrics      []InstanceMetrics
}

type InstanceMetrics struct {
	Name string
	Port int64
	Path string
}

type DiscoveredInstances struct {
	Instances []Instance
}

type outputFormat struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func ToJsonString(d DiscoveredInstances) ([]byte, error) {
	outputList := []outputFormat{}

	for _, instance := range d.Instances {
		for _, metric := range instance.Metrics {
			targetHost := fmt.Sprintf("%s:%d", instance.PrivateIP, metric.Port)
			l := labels(instance.Tags, targetHost, metric.Path, metric.Name)

			output := outputFormat{Targets: []string{targetHost}, Labels: l}
			outputList = append(outputList, output)
		}
	}

	return json.Marshal(outputList)
}

func labels(tags map[string]string, targetHost string, path string, metricName string) map[string]string {
	labels := lowerMapKeys(tags)
	labels["__metrics_path__"] = path
	labels["__address__"] = targetHost
	if val, ok := labels["name"]; ok {
		labels["instancename"] = val
	}
	labels["name"] = metricName
	return labels
}

func lowerMapKeys(tags map[string]string) map[string]string {
	lowerMap := map[string]string{}
	for key, value := range tags {
		lowerMap[strings.ToLower(key)] = value
	}
	return lowerMap
}
