package outputkubernetes

import (
	"context"
	"testing"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/endpoints"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/libtesting"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	fake "k8s.io/client-go/kubernetes/fake"
)

func TestWriteShouldCreateConfigmap(t *testing.T) {
	var clientset kubernetes.Interface
	clientset = fake.NewSimpleClientset()

	o := OutputKubernetes{
		clientset:      clientset,
		ConfigMapField: "discovery",
		ConfigMapName:  "testcm",
		Namespace:      "default",
	}

	iList := endpoints.DiscoveredInstances{
		Instances: libtesting.InstanceList(),
	}
	o.Write(iList)

	resp, _ := clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), "testcm", metav1.GetOptions{})
	if resp == nil {
		t.Error("Expected to create configmap with name 'testcm' in namespace 'default', but no found")
	}

	if _, ok := resp.Data["discovery"]; !ok {
		t.Error("Expected to see field 'discovery' in configmap to be defined")
	}
}

func TestWriteShouldUpdateConfigmap(t *testing.T) {
	cmData := make(map[string]string)
	cmData["key"] = "value"

	var clientset kubernetes.Interface
	clientset = fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testcm",
			Namespace: "default",
		},
		Data: cmData,
	})

	o := OutputKubernetes{
		clientset:      clientset,
		ConfigMapField: "discovery",
		ConfigMapName:  "testcm",
		Namespace:      "default",
	}
	iList := endpoints.DiscoveredInstances{
		Instances: libtesting.InstanceList(),
	}
	o.Write(iList)

	resp, _ := clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), "testcm", metav1.GetOptions{})

	if resp.Data["key"] != "value" {
		t.Error("Expected to see an update of an existing configmap that adds only new field, but existing field seems to be overwritten")
	}

	if _, ok := resp.Data["discovery"]; !ok {
		t.Error("Expected to see field 'discovery' in configmap to be defined")
	}
}
