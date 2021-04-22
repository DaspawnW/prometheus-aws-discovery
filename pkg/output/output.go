package output

import "github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"

type Output interface {
	Write([]discovery.Instance) error
}
