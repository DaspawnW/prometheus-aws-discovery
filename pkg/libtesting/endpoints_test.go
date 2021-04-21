package libtesting

import (
	"encoding/json"
	"testing"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/endpoints"
	"github.com/nsf/jsondiff"
)

func TestToJsonStringHasCorrectContent(t *testing.T) {

	instanceList := InstanceList()

	returnedJSONContentBytes, _ := endpoints.ToJSONString(endpoints.DiscoveredInstances{Instances: instanceList})
	expectedJson := []map[string]interface{}{
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
			},
		},
		{
			"targets": []string{
				"127.0.0.1:8080",
			},
			"labels": map[string]string{
				"__address__":      "127.0.0.1:8080",
				"__metrics_path__": "/metrics",
				"__scheme__":       "https",
				"billingnumber":    "1111",
				"instancename":     "Testinstance1",
				"name":             "blackbox_exporter",
			},
		},
		{
			"targets": []string{
				"127.0.0.2:9100",
			},
			"labels": map[string]string{
				"__address__":      "127.0.0.2:9100",
				"__metrics_path__": "/metrics",
				"__scheme__":       "http",
				"billingnumber":    "2222",
				"instancename":     "Testinstance2",
				"name":             "node_exporter",
			},
		},
	}
	expectedJsonBytes, _ := json.Marshal(expectedJson)
	opts := jsondiff.DefaultConsoleOptions()
	res, text := jsondiff.Compare(expectedJsonBytes, returnedJSONContentBytes, &opts)
	if res != 0 {
		t.Errorf("JsonDiff \n%v \n%v", res, text)
	}
	badJson := []map[string]interface{}{
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
	badJsonBytes, _ := json.Marshal(badJson)
	res, text = jsondiff.Compare(badJsonBytes, returnedJSONContentBytes, &opts)
	if res == 0 {
		t.Errorf("JsonDiff \n%v \n%v", res, text)
	}

}
