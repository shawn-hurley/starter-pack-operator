package broker

import (
	"fmt"
	"reflect"

	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	"github.com/operator-framework/operator-sdk/pkg/sdk/query"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	extapi "k8s.io/api/extensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
	}
	// create the spec
	spec := extapi.DeploymentSpec{
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
						VolumeMounts: createVolumeMounts(br),
					},
				},
				Volumes: createVolumes(br),
			},
		},
	}
	err := query.Get(d)
	if apierrors.IsNotFound(err) {
		d.Spec = spec
		addOwnerRefToObject(d, asOwner(br))
		err := action.Create(d)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			logrus.Debugf("Deployment unable to be created: %v", err)
			return err
		}
		return nil
	}
	if err != nil {
		logrus.Debugf("Deployment unable to be created: %v", err)
		return err
	}
	if !reflect.DeepEqual(d.Spec, spec) {
		//If the specs are not equal then something in broker spec has changed.
		// We should update the deployment then.
		d.Spec = spec
		err := action.Update(d)
		if err != nil {
			logrus.Debugf("Deployment unable to be created: %v", err)
			return err
		}
		return nil
	}
	return nil
}

func createArgs(br *api.Broker) []string {
	args := []string{"--port", fmt.Sprintf("%v", br.Spec.Port)}
	if br.Spec.AuthenticateK8SToken {
		args = append(args, "--authenticate-k8s-token")
	}
	if br.Spec.TLSSecretRef != nil {
		args = append(args, "--tls-cert-file", "/var/run/osb-starter-pack/cert.crt", "--tls-private-key-file", "/var/run/osb-starter-pack/cert.key")
	}
	args = append(args, "-v", "5", "-logtostderr")
	return args
}

func createVolumeMounts(br *api.Broker) []v1.VolumeMount {
	if br.Spec.TLSSecretRef != nil {
		return []v1.VolumeMount{
			v1.VolumeMount{
				MountPath: "/var/run/osb-starter-pack",
				Name:      "osb-starter-pack-ssl",
				ReadOnly:  true,
			},
		}
	}
	return nil
}

func createVolumes(br *api.Broker) []v1.Volume {
	if br.Spec.TLSSecretRef != nil {
		var mode int32 = 420
		return []v1.Volume{
			v1.Volume{
				Name: "osb-starter-pack-ssl",
				VolumeSource: v1.VolumeSource{
					Secret: &v1.SecretVolumeSource{
						DefaultMode: &mode,
						SecretName:  fmt.Sprintf("tls-%v", br.Name),
						Items: []v1.KeyToPath{
							v1.KeyToPath{
								Key:  "cert.crt",
								Path: "cert.crt",
							},
							v1.KeyToPath{
								Key:  "cert.key",
								Path: "cert.key",
							},
						},
					},
				},
			},
		}
	}

	return nil
}
