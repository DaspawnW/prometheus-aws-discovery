package output

import (
	"fmt"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/endpoints"
	log "github.com/sirupsen/logrus"
)

func Write(instances endpoints.DiscoveredInstances) error {
	content, err := endpoints.ToJSONString(instances)
	if err != nil {
		log.Error("Failed to convert instances to json string with error", err)
		return err
	}
	fmt.Printf("%v", content)
	return nil
}
