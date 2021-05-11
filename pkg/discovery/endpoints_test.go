package discovery

import (
	"encoding/json"
	"testing"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/test"
	"github.com/mitchellh/mapstructure"
	"github.com/nsf/jsondiff"
)

func TestJSONStringHasCorrectContent(t *testing.T) {
	instances := make([]Instance, 0)
	for _, instanceData := range test.Instances {
		var instance Instance
		mapstructure.Decode(instanceData, &instance)
		instances = append(instances, instance)
	}
	returnedJSONContentBytes, _ := TargetConfigBytes(instances)
	expectedJSON := test.Targets
	expectedJSONBytes, _ := json.Marshal(expectedJSON)
	opts := jsondiff.DefaultConsoleOptions()
	res, text := jsondiff.Compare(returnedJSONContentBytes, expectedJSONBytes, &opts)
	if res != 0 {
		t.Errorf("JSONDiff \n%v \n%v", res, text)
	}

	badJSON := []map[string]interface{}{
		{
			"targets": []string{
				"127.0.0.1:9100",
			},
			"labels": map[string]string{
				"__address__":      "127.0.0.1:9100",
				"__metrics_path__": "/metrics",
				"__scheme__":       "http",
				"billingnumber":    "1111",
				"instancename":     "Testinstance1",
				"name":             "node_exporter",
				"spotprice":        "123",
			},
		},
	}
	badJSONBytes, _ := json.Marshal(badJSON)
	res, text = jsondiff.Compare(badJSONBytes, returnedJSONContentBytes, &opts)
	if res == 0 {
		t.Errorf("JSONDiff \n%v \n%v", res, text)
	}

}
