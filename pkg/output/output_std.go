package output

import (
	"fmt"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"
)

type OutputStdOut struct {
}

func (o OutputStdOut) Write(instances []discovery.Instance) error {
	jsonBytes, err := discovery.TargetConfigBytes(instances)
	if err != nil {
		return err
	}
	fmt.Printf("Discovered Scrape Configs: \n%v\n", string(jsonBytes))
	return nil
}
