package resource

import (
	"context"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	e "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesResourceSecret struct {
	clientset kubernetes.Interface
}

func NewKubernetesResourceSecret(clientset kubernetes.Interface) KubernetesResource {
	return KubernetesResourceSecret{
		clientset: clientset,
	}
}

func (k KubernetesResourceSecret) Create(namespace string, name string, field string, data string) error {
	cmData := make(map[string][]byte)
	cmData[field] = []byte(data)
	cm := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: cmData,
	}

	log.Debugf("Create secret with name %s in namespace %s", cm.ObjectMeta.Name, cm.ObjectMeta.Namespace)
	_, err := k.clientset.CoreV1().Secrets(namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
	return err
}

func (k KubernetesResourceSecret) Update(namespace string, name string, field string, data string) error {
	cm, err := k.clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if cm.Data == nil {
		cm.Data = make(map[string][]byte)
	}

	cm.Data[field] = []byte(data)

	log.Debugf("Update secret %s in namespace %s", cm.ObjectMeta.Namespace, cm.ObjectMeta.Name)
	_, saveErr := k.clientset.CoreV1().Secrets(cm.ObjectMeta.Namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	return saveErr
}

func (k KubernetesResourceSecret) Exists(namespace string, name string) (bool, error) {
	_, err := k.clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if !e.IsNotFound(err) {
			log.Error("Failed to load secret from k8s", err)
			return false, err
		}

		return false, nil
	}

	return true, nil
}
