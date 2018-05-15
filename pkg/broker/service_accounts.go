package broker

import (
	"fmt"

	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	"k8s.io/api/core/v1"
	authz "k8s.io/api/rbac/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func syncClientServiceAccount(br *api.Broker) error {
	// Create client service account
	sc := &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-client", br.Name),
			Namespace: br.Namespace,
		},
	}
	clusterRole := &authz.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authroization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-client", br.Name),
			Namespace: br.Namespace,
		},
		Rules: []authz.PolicyRule{
			authz.PolicyRule{
				NonResourceURLs: []string{"/v2", "/v2/*"},
				Verbs:           []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			},
		},
	}
	clusterRoleBinding := &authz.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authroization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-client", br.Name),
			Namespace: br.Namespace,
		},
		Subjects: []authz.Subject{
			authz.Subject{
				Kind:      "ServiceAccount",
				Name:      fmt.Sprintf("%v-client", br.Name),
				Namespace: br.Namespace,
			},
		},
		RoleRef: authz.RoleRef{
			APIGroup: "rbac.Authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     fmt.Sprintf("%v-client", br.Name),
		},
	}
	scSecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-client-secret", br.Name),
			Namespace: br.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": fmt.Sprintf("%v-client", br.Name),
			},
		},
		Type: v1.SecretTypeServiceAccountToken,
	}
	err := action.Create(sc)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("Unable to create brokers client service account: %v", err)

	}
	err = action.Create(clusterRole)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("Unable to create brokers client cluser role: %v", err)
	}
	err = action.Create(clusterRoleBinding)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("Unable to create brokers client cluster role binding: %v", err)
	}
	err = action.Create(scSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("Unable to create brokers client secret: %v", err)
	}
	return nil
}

func syncBrokerServiceAccount(br *api.Broker) error {
	sc := &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-service", br.Name),
			Namespace: br.Namespace,
		},
	}

	err := action.Create(sc)
	if err != nil {
		return fmt.Errorf("Unable to create brokers service service account: %v", err)
	}

	// If we are not authenticationg the token then we don't need any special parameters.
	if br.Spec.AuthenticateK8SToken {
		clusterRole := &authz.ClusterRole{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterRole",
				APIVersion: "rbac.authroization.k8s.io/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%v-service", br.Name),
				Namespace: br.Namespace,
			},
			Rules: []authz.PolicyRule{
				authz.PolicyRule{
					APIGroups: []string{"authentication.k8s.io"},
					Resources: []string{"tokenreviews"},
					Verbs:     []string{"create"},
				},
				authz.PolicyRule{
					APIGroups: []string{"authorization.k8s.io"},
					Resources: []string{"subjectaccessreviews"},
					Verbs:     []string{"create"},
				},
			},
		}
		err := action.Create(clusterRole)
		if err != nil {
			return fmt.Errorf("Unable to create brokers service cluster role: %v", err)
		}

		clusterRoleBinding := &authz.ClusterRoleBinding{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterRoleBinding",
				APIVersion: "rbac.authroization.k8s.io/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%v-service", br.Name),
				Namespace: br.Namespace,
			},
			Subjects: []authz.Subject{
				authz.Subject{
					Kind:      "ServiceAccount",
					Name:      fmt.Sprintf("%v-service", br.Name),
					Namespace: br.Namespace,
				},
			},
			RoleRef: authz.RoleRef{
				APIGroup: "rbac.Authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     fmt.Sprintf("%v-service", br.Name),
			},
		}
		err = action.Create(clusterRoleBinding)
		if err != nil {
			return fmt.Errorf("Unable to create brokers service cluster role binding: %v", err)
		}
	}
	return nil
}
