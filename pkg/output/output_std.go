package output

import (
	"fmt"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"
	log "github.com/sirupsen/logrus"
)

func Write(instances []discovery.Instance) error {
	content, err := discovery.TargetConfigBytes(instances)
	if err != nil {
		log.Error("Failed to convert instances to json string with error", err)
		return err
	}
	fmt.Printf("%v", content)
	return nil
}
