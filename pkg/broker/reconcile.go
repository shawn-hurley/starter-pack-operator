package broker

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
)

// Reconcile - reconciles the brokers state accross the cluster service broker
// and the broker pods service account.
func Reconcile(br *api.Broker) error {
	br = br.DeepCopy()
	changed := br.SetDefaults()
	if changed {
		return action.Update(br)
	}
	// After first time reconcile, phase will switch to Running
	if br.Status.Phase == api.PhaseInitial {
		// TODO: if we should generate a cert, shoud generate a ca, cert, key
		//This should update the defaults. Also should update the status to creating.
		prepareBrokerTLSSecrets(br)
	}

	//Determine if we need to create the client service account
	if br.Spec.AuthenticateK8SToken {
		err := syncClientServiceAccount(br)
		if err != nil {
			return err
		}
	}

	err := syncBrokerServiceAccount(br)
	if err != nil {
		return err
	}

	err = syncBrokerDeployment(br)
	if err != nil {
		return err
	}

	err = syncBrokerService(br)
	if err != nil {
		return err
	}

	err = syncClusterServiceBroker(br)
	if err != nil {
		return err
	}
	return nil
}
