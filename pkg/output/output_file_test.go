package output

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"
)

func TestWriteToFile(t *testing.T) {

	instances := []discovery.Instance{
		{
			InstanceType: "t2.small",
			PrivateIP:    "127.0.0.2",
			Tags: map[string]string{
				"test": "tags",
			},
			Metrics: []discovery.InstanceMetrics{
				{
					Name:   "someName",
					Port:   8888,
					Path:   "/metrics",
					Scheme: "https",
				},
			},
		},
	}

	path := "_test.json"
	o := OutputFile{
		FilePath: path,
	}

	// initial cleanup
	os.Remove(path)

	o.Write(instances)

	bs, err := ioutil.ReadFile(path)

	if err != nil || len(bs) == 0 {
		t.Error("Failed to write instances list to file", err)
	}

	os.Remove(path)
	// final cleanup
}
