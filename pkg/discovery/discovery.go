package discovery

type DiscoveryClient interface {
	GetInstances() ([]Instance, error)
}
