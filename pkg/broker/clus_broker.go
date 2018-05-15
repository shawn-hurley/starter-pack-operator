package broker

import (
	"fmt"

	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	"k8s.io/api/core/v1"
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
	caBundle := fmt.Sprintf("%s", se.Data["ca.crt"])

	u := unstructured.Unstructured{}
	u.SetAPIVersion("servicecatalog.k8s.io/v1beta1")
	u.SetKind("ClusterServiceBroker")
	u.SetName(br.Name)
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
	return action.Create(&u)
}
