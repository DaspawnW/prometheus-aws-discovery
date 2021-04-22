package discovery

type DiscoveryClient interface {
	newDiscovery(client interface{}, prefix string) DiscoveryClient
	getTargets() []Instance
}
