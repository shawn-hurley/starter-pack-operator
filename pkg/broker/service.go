package broker

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func syncBrokerService(br *api.Broker) error {
	s := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      br.Name,
			Namespace: br.Namespace,
			Labels: map[string]string{
				"app": br.Name,
			},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app": br.Name,
			},
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Protocol: v1.ProtocolTCP,
					Port:     int32(443),
					TargetPort: intstr.IntOrString{
						IntVal: int32(br.Spec.Port),
					},
				},
			},
		},
	}
	return action.Create(s)
}
