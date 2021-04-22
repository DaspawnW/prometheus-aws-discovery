package discovery

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
	Name   string
	Port   int64
	Path   string
	Scheme string
}

type outputFormat struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func TargetConfigBytes(d []Instance) ([]byte, error) {
	outputList := []outputFormat{}

	for _, instance := range d {
		for _, metric := range instance.Metrics {
			targetHost := fmt.Sprintf("%s:%d", instance.PrivateIP, metric.Port)
			l := labels(instance.Tags, targetHost, metric.Path, metric.Name, metric.Scheme)

			output := outputFormat{Targets: []string{targetHost}, Labels: l}
			outputList = append(outputList, output)
		}
	}

	return json.Marshal(outputList)
}

func labels(tags map[string]string, targetHost string, path string, metricName string, scheme string) map[string]string {
	labels := standardizeKeys(tags)
	labels["__metrics_path__"] = path
	labels["__scheme__"] = scheme
	labels["__address__"] = targetHost
	if val, ok := labels["name"]; ok {
		labels["instancename"] = val
	}
	labels["name"] = metricName
	return labels
}

func standardizeKeys(tags map[string]string) map[string]string {
	lowerMap := map[string]string{}
	for key, value := range tags {
		reducedByEmpty := strings.ReplaceAll(key, " ", "")
		if len(reducedByEmpty) > 0 {
			lowerMap[strings.ToLower(reducedByEmpty)] = value
		}
	}
	return lowerMap
}
