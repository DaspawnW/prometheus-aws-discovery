package output

import "github.com/daspawnw/prometheus-aws-discovery/pkg/endpoints"

type Output interface {
	Write(endpoints.DiscoveredInstances) error
}
