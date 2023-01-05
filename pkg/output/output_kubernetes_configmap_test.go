package output

import (
	"context"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/output/resource"
	"testing"

	"github.com/daspawnw/prometheus-aws-discovery/pkg/discovery"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	fake "k8s.io/client-go/kubernetes/fake"
)

var iList = []discovery.Instance{
	{
		InstanceType: "t2.small",
		PrivateIP:    "127.0.0.2",
		Tags: map[string]string{
			"test": "tags",
		},
		Metrics: []discovery.InstanceMetrics{
			{
				Name:   "someName",
				Port:   8888,
				Path:   "/metrics",
				Scheme: "https",
			},
		},
	},
}

func TestWriteShouldCreateConfigmap(t *testing.T) {
	var clientset kubernetes.Interface
	clientset = fake.NewSimpleClientset()

	r := resource.NewKubernetesResourceConfigmap(clientset)
	o := OutputKubernetes{
		resourceOperation: r,
		ResourceField:     "discovery",
		ResourceName:      "testcm",
		Namespace:         "default",
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

	r := resource.NewKubernetesResourceConfigmap(clientset)
	o := OutputKubernetes{
		resourceOperation: r,
		ResourceField:     "discovery",
		ResourceName:      "testcm",
		Namespace:         "default",
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
