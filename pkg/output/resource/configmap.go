package resource

import (
	"context"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	e "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesResourceConfigmap struct {
	clientset kubernetes.Interface
}

func NewKubernetesResourceConfigmap(clientset kubernetes.Interface) KubernetesResource {
	return KubernetesResourceConfigmap{
		clientset: clientset,
	}
}

func (k KubernetesResourceConfigmap) Create(namespace string, name string, field string, data string) error {
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
	_, err := k.clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
	return err
}

func (k KubernetesResourceConfigmap) Update(namespace string, name string, field string, data string) error {
	cm, err := k.clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	cm.Data[field] = data

	log.Debugf("Update configmap %s in namespace %s", cm.ObjectMeta.Namespace, cm.ObjectMeta.Name)
	_, saveErr := k.clientset.CoreV1().ConfigMaps(cm.ObjectMeta.Namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	return saveErr
}

func (k KubernetesResourceConfigmap) Exists(namespace string, name string) (bool, error) {
	_, err := k.clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if !e.IsNotFound(err) {
			log.Error("Failed to load configmap from k8s", err)
			return false, err
		}

		return false, nil
	}

	return true, nil
}
