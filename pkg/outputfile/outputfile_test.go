package outputfile

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/endpoints"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/libtesting"
)

func TestWriteToFile(t *testing.T) {
	instances := endpoints.DiscoveredInstances{
		Instances: libtesting.InstanceList(),
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
