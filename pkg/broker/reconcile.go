package broker

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk/action"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	log "github.com/sirupsen/logrus"
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
		log.Infof("initial phase.")
		// TODO: if we should generate a cert, shoud generate a ca, cert, key
		//This should update the defaults. Also should update the status to creating.
		err := prepareBrokerTLSSecrets(br)
		if err != nil {
			log.Errorf("unable to prepare borker TLS Secret")
			br.Status.Phase = api.PhaseError
			action.Update(br)
			return err
		}
		br.Status.Phase = api.PhaseCreating
		action.Update(br)
		return nil
	}

	log.Infof("create broker phase.")

	//Determine if we need to create the client service account
	if br.Spec.AuthenticateK8SToken {
		err := syncClientServiceAccount(br)
		if err != nil {
			br.Status.Phase = api.PhaseError
			action.Update(br)
			return err
		}
	}

	log.Infof("sync service account.")
	err := syncBrokerServiceAccount(br)
	if err != nil {
		br.Status.Phase = api.PhaseError
		action.Update(br)
		return err
	}

	log.Infof("sync deployment.")
	err = syncBrokerDeployment(br)
	if err != nil {
		br.Status.Phase = api.PhaseError
		action.Update(br)
		return err
	}

	log.Infof("sync broker service.")
	err = syncBrokerService(br)
	if err != nil {
		br.Status.Phase = api.PhaseError
		action.Update(br)
		return err
	}

	log.Infof("sync cluster service broker resource.")
	err = syncClusterServiceBroker(br)
	if err != nil {
		br.Status.Phase = api.PhaseError
		action.Update(br)
		return err
	}
	// Don't update resource if nothing needs to change.
	if br.Status.Phase != api.PhaseRunning {
		br.Status.Phase = api.PhaseRunning
		action.Update(br)
	}
	return nil
}
