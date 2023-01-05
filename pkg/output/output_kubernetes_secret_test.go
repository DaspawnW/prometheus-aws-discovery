package output

import (
	"context"
	"encoding/base64"
	"github.com/daspawnw/prometheus-aws-discovery/pkg/output/resource"
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	fake "k8s.io/client-go/kubernetes/fake"
)

func TestWriteShouldCreateSecret(t *testing.T) {
	var clientset kubernetes.Interface
	clientset = fake.NewSimpleClientset()

	r := resource.NewKubernetesResourceSecret(clientset)
	o := OutputKubernetes{
		resourceOperation: r,
		ResourceField:     "discovery",
		ResourceName:      "testsecret-new",
		Namespace:         "default",
	}

	o.Write(iList)

	resp, _ := clientset.CoreV1().Secrets("default").Get(context.TODO(), "testsecret-new", metav1.GetOptions{})
	if resp == nil {
		t.Error("Expected to create configmap with name 'testsecret-new' in namespace 'default', but no found")
	}

	if _, ok := resp.Data["discovery"]; !ok {
		t.Error("Expected to see field 'discovery' in configmap to be defined")
	}
}

func TestWriteShouldUpdateSecret(t *testing.T) {
	base64Encoded := make([]byte, base64.StdEncoding.EncodedLen(len("value")))
	base64.StdEncoding.Encode(base64Encoded, []byte("value"))

	cmData := make(map[string][]byte)
	cmData["key"] = base64Encoded

	var clientset kubernetes.Interface
	clientset = fake.NewSimpleClientset(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testsecret-exists",
			Namespace: "default",
		},
		Data: cmData,
	})

	r := resource.NewKubernetesResourceSecret(clientset)
	o := OutputKubernetes{
		resourceOperation: r,
		ResourceField:     "discovery",
		ResourceName:      "testsecret-exists",
		Namespace:         "default",
	}

	o.Write(iList)

	resp, _ := clientset.CoreV1().Secrets("default").Get(context.TODO(), "testsecret-exists", metav1.GetOptions{})

	base64Decoded := make([]byte, base64.StdEncoding.DecodedLen(len(resp.Data["key"])))
	base64.StdEncoding.Decode(base64Decoded, resp.Data["key"])
	s := string(base64Decoded)

	if !strings.Contains(s, "value") {
		t.Error("Expected to see an update of an existing secret that adds only new field, but existing field seems to be overwritten")
	}

	if _, ok := resp.Data["discovery"]; !ok {
		t.Error("Expected to see field 'discovery' in configmap to be defined")
	}
}
