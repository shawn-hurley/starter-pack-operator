package broker

import (
	"reflect"

	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
	}
	spec := v1.ServiceSpec{
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
	}

	err := query.Get(s)
	if apierrors.IsNotFound(err) {
		s.Spec = spec
		addOwnerRefToObject(s, asOwner(br))
		err := action.Create(s)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			log.Errorf("unable to create service - %v", err)
			return err
		}
	}

	if !reflect.DeepEqual(s.Spec.Ports, spec.Ports) {
		s.Spec.Ports = spec.Ports
		err := action.Update(s)
		if err != nil {
			log.Errorf("unable to update service - %v", err)
			return err
		}
		return nil
	}
	return nil
}
