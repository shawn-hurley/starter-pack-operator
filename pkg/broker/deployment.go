package broker

import (
	"fmt"

	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	"k8s.io/api/core/v1"
	extapi "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func syncBrokerDeployment(br *api.Broker) error {
	var replicas int32
	replicas = 1
	d := &extapi.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      br.Name,
			Namespace: br.Namespace,
		},
		Spec: extapi.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": br.Name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": br.Name,
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: fmt.Sprintf("%v-service", br.Name),
					Containers: []v1.Container{
						v1.Container{
							Name:            "service-broker-skeleton",
							Image:           br.Spec.Image,
							ImagePullPolicy: v1.PullAlways,
							Command:         []string{"/opt/servicebroker/servicebroker"},
							Args:            createArgs(br),
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									ContainerPort: int32(br.Spec.Port),
								},
							},
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									TCPSocket: &v1.TCPSocketAction{
										Port: intstr.IntOrString{
											IntVal: int32(br.Spec.Port),
										},
									},
								},
								FailureThreshold:    int32(1),
								InitialDelaySeconds: int32(10),
								PeriodSeconds:       int32(10),
								SuccessThreshold:    int32(1),
								TimeoutSeconds:      int32(2),
							},
						},
					},
				},
			},
		},
	}
	return action.Create(d)
	return nil
}

func createArgs(br *api.Broker) []string {
	args := []string{"--port", fmt.Sprintf("%v", br.Spec.Port)}
	if br.Spec.TLSCert != "" {
		args = append(args, "--tlsCert", br.Spec.TLSCert)
	}
	if br.Spec.TLSKey != "" {
		args = append(args, "--tlsKey", br.Spec.TLSKey)
	}
	if br.Spec.AuthenticateK8SToken {
		args = append(args, "--authenticate-k8s-token")
	}
	if br.Spec.TLSSecretRef != nil {
		args = append(args, "--tls-cert-file", "/var/run/osb-starter-pack/tls.crt", "--tls-private-key-file", "/var/run/osb-starter-pack/tls.key")
	}
	args = append(args, "-v", "5", "-logtostderr")
	return args
}
