package output

import (
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

func getClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		log.Debug("Use kubeconfig provided by commandline flag")
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	log.Debug("Use in-cluster k8s configuration")
	return rest.InClusterConfig()
}

func getNamespaceByServiceAccount() string {
	nsBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Debugf("No Namespace provided and no namespace found in-cluster path, using %s", apiv1.NamespaceDefault)
		return apiv1.NamespaceDefault
	}

	log.Debug("Detected namespace in-cluster path")
	return strings.TrimSpace(string(nsBytes))
}
