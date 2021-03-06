package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrokerList - list of the brokers
type BrokerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Broker `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Broker - list of the brokers
type Broker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              BrokerSpec   `json:"spec"`
	Status            BrokerStatus `json:"status,omitempty"`
}

// SetDefaults - Set the defaults for the broker.
func (b *Broker) SetDefaults() bool {
	changed := false
	bs := &b.Spec
	if bs.Port == 0 {
		bs.Port = 1338
		changed = true
	}
	if bs.Image == "" {
		bs.Image = "quay.io/osb-starter-pack/servicebroker:latest"
		changed = true
	}
	return changed
}

// BrokerSpec - Used to specify the deployment of a broker.
type BrokerSpec struct {
	// If you wish to specify the port that the broker will listen on.
	// Defaults to 1338
	Port int `json:"port"`
	// Reference to a secret with a tls.crt and tls.key
	// Should also contain a ca.crt.
	TLSSecretRef         *v1.ObjectReference `json:"tlsSecretRef"`
	Image                string              `json:"image"`
	AuthenticateK8SToken bool                `json:"authenticateK8SToken"`
}

// BrokerPhase - the status phase of the broker
type BrokerPhase string

const (
	// PhaseInitial - broker phase initial.
	PhaseInitial BrokerPhase = ""
	// PhaseCreating - broker phase while broker is being created.
	PhaseCreating BrokerPhase = "Creating"
	// PhaseRunning - broker phase describes when a broker is running.
	// Will be moved when the status of svc cat clusterservicebroker
	// Is retrieving values.
	PhaseRunning BrokerPhase = "Running"
	// PhaseError - broker error describes that the broker is in an error
	// State. This could mean the pod is not comming up or the
	// clusterservicebroker is not able to contanct the broker.
	PhaseError BrokerPhase = "Error"
)

// BrokerStatus - broker status
type BrokerStatus struct {
	Phase BrokerPhase `json:"phase"`
}
