package output

import (
	"errors"
	"fmt"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/output/resource"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type OutputKubernetes struct {
	resourceOperation resource.KubernetesResource
	Namespace         string
	ResourceName      string
	ResourceField     string
}

type OutputKubernetesResourceType string

const (
	KUBERNETES_CONFIGMAP OutputKubernetesResourceType = "configmap"
	KUBERNETES_SECRET    OutputKubernetesResourceType = "secret"
)

func NewOutputKubernetes(kubeconfig string, resourceType OutputKubernetesResourceType, namespace string, resourceName string, resourceField string) (*OutputKubernetes, error) {
	config, err := getClientConfig(kubeconfig)
	if err != nil {
		log.Error("Failed to configure k8s client")
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("Failed to create k8s clientset", err)
		return nil, err
	}

	var r resource.KubernetesResource
	if resourceType == KUBERNETES_CONFIGMAP {
		r = resource.NewKubernetesResourceConfigmap(clientset)
	} else if resourceType == KUBERNETES_SECRET {
		r = resource.NewKubernetesResourceSecret(clientset)
	} else {
		log.Error("Invalid Kubernetes output format provided")
		return nil, errors.New(fmt.Sprintf("Invalid Kubernetes output format provided %s", resourceType))
	}

	return &OutputKubernetes{
		resourceOperation: r,
		ResourceField:     resourceField,
		ResourceName:      resourceName,
		Namespace:         namespace,
	}, nil
}

func (o OutputKubernetes) Write(instances []discovery.Instance) error {
	output, err := discovery.TargetConfigBytes(instances)
	if err != nil {
		log.Error("Failed to convert instances to json string")
		return err
	}

	ns := o.getNamespace()

	exists, err := o.resourceOperation.Exists(ns, o.ResourceName)
	if err != nil {
		log.Error("Failed to load from k8s", err)
		return err
	}

	if exists == true {
		err := o.resourceOperation.Update(ns, o.ResourceName, o.ResourceField, string(output))
		if err != nil {
			log.Error("Failed to update resource")
		}
		return err
	} else {
		err := o.resourceOperation.Create(ns, o.ResourceName, o.ResourceField, string(output))
		if err != nil {
			log.Error("Failed to create resource")
		}
		return err
	}
}

func (o OutputKubernetes) getNamespace() string {
	if o.Namespace != "" {
		log.Debug("Use Namespace provided by commandline flag")
		return o.Namespace
	}

	return getNamespaceByServiceAccount()
}
