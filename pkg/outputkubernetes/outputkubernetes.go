package outputkubernetes

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/endpoints"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type OutputKubernetes struct {
	clientset      kubernetes.Interface
	Namespace      string
	ConfigMapName  string
	ConfigMapField string
}

func NewOutputKubernetes(kubeconfig string, namespace string, cmName string, cmField string) (*OutputKubernetes, error) {
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

	return &OutputKubernetes{
		clientset:      clientset,
		ConfigMapField: cmField,
		ConfigMapName:  cmName,
		Namespace:      namespace,
	}, nil
}

func (o OutputKubernetes) Write(instances endpoints.DiscoveredInstances) error {
	output, err := endpoints.ToJsonString(instances)
	if err != nil {
		log.Error("Failed to convert instances to json string")
		return err
	}

	ns := o.getNamespace()
	cm, err := loadConfigmap(o.clientset, ns, o.ConfigMapName)
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Error("Failed to load configmap from k8s", err)
			return err
		}
		log.Infof("No configmap with name %s in namespace %s detected", o.ConfigMapName, ns)
	}

	if cm != nil {
		err := updateConfigmap(o.clientset, cm, o.ConfigMapField, string(output))
		if err != nil {
			log.Error("Failed to update configmap")
		}
		return err
	} else {
		err := createConfigmap(o.clientset, o.ConfigMapName, ns, o.ConfigMapField, string(output))
		if err != nil {
			log.Error("Failed to create configmap")
		}
		return err
	}
}

func (o OutputKubernetes) getNamespace() string {
	if o.Namespace != "" {
		log.Debug("Use Namespace provided by commandline flag")
		return o.Namespace
	}

	nsBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Debugf("No Namespace provided and no namespace found in-cluster path, using %s", apiv1.NamespaceDefault)
		return apiv1.NamespaceDefault
	}

	log.Debug("Detected namespace in-cluster path")
	return strings.TrimSpace(string(nsBytes))
}

func createConfigmap(clientset kubernetes.Interface, name string, namespace string, field string, data string) error {
	cmData := make(map[string]string)
	cmData[field] = data
	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: cmData,
	}
	log.Debugf("Create configmap with name %s in namespace %s", cm.ObjectMeta.Name, cm.ObjectMeta.Namespace)
	_, err := clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
	return err
}

func updateConfigmap(clientset kubernetes.Interface, cm *v1.ConfigMap, field string, data string) error {
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}
	cm.Data[field] = data
	log.Debugf("Update configmap %s in namespace %s", cm.ObjectMeta.Namespace, cm.ObjectMeta.Name)
	_, err := clientset.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	return err
}

func loadConfigmap(clientset kubernetes.Interface, namespace string, name string) (*v1.ConfigMap, error) {
	cm, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func getClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		log.Debug("Use kubeconfig provided by commandline flag")
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	log.Debug("Use in-cluster k8s configuration")
	return rest.InClusterConfig()
}
