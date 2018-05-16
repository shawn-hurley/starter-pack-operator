package broker

import (
	"encoding/base64"
	"fmt"

	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func syncClusterServiceBroker(br *api.Broker) error {
	// Get the CA data from the secret.
	se := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("tls-%v", br.Name),
			Namespace: br.Namespace,
		},
	}
	err := query.Get(se)
	if err != nil {
		return err
	}
	ca := se.Data["ca.crt"]
	logrus.Infof("ca: %v", ca)
	logrus.Infof("ca string: %q", ca)
	caBundle := base64.StdEncoding.EncodeToString(se.Data["ca.crt"])
	logrus.Infof("ca bundle: %v", caBundle)

	u := &unstructured.Unstructured{}
	spec := map[string]interface{}{
		"url":      fmt.Sprintf("https://%v.%v.svc.cluster.local", br.Name, br.Namespace),
		"caBundle": caBundle,
	}

	if br.Spec.AuthenticateK8SToken {
		spec["authInfo"] = map[string]interface{}{
			"bearer": map[string]interface{}{
				"secretRef": map[string]interface{}{
					"namespace": br.Namespace,
					"name":      fmt.Sprintf("%v-client-secret", br.Name),
				},
			},
		}
	}

	c := map[string]interface{}{
		"spec": spec,
	}
	u.SetUnstructuredContent(c)
	u.SetAPIVersion("servicecatalog.k8s.io/v1beta1")
	u.SetKind("ClusterServiceBroker")
	u.SetName(br.Name)
	logrus.Infof("cluster service broker: %#v", u)
	err = action.Create(u)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		logrus.Errorf("unable to create cluster service broker: %v", err)
		return err
	}
	return nil
}
