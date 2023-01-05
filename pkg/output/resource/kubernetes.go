package resource

type KubernetesResource interface {
	Exists(namespace string, name string) (bool, error)
	Create(namespace string, name string, field string, data string) error
	Update(namespace string, name string, field string, data string) error
}
