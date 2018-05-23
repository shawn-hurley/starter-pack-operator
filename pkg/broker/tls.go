package broker

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	action "github.com/operator-framework/operator-sdk/pkg/sdk"
	api "github.com/shawn-hurley/starter-pack-operator/pkg/apis/starterpack/v1alpha1"
	"github.com/shawn-hurley/starter-pack-operator/pkg/tls"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func prepareBrokerTLSSecrets(br *api.Broker) error {
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

	err := action.Get(se)
	if err == nil {
		return nil
	}
	if !apierrors.IsNotFound(err) {
		return err
	}

	caKey, caCrt, err := newCACert(br)
	tc := tls.CertConfig{
		CommonName:   fmt.Sprintf("%v operator", br.Name),
		Organization: []string{"starterpack.org"},
		AltNames:     tls.NewAltNames([]string{fmt.Sprintf("*.%v.%v.svc.cluster.local", br.Name, br.Namespace), fmt.Sprintf("%v.%v.svc.cluster.local", br.Name, br.Namespace)}),
	}
	key, crt, err := newKeyAndCert(caCrt, caKey, tc)
	if err != nil {
		return fmt.Errorf("new TLS secret failed: %v", err)
	}

	se.Data = map[string][]byte{
		"cert.key": tls.EncodePrivateKeyPEM(key),
		"cert.crt": tls.EncodeCertificatePEM(crt),
		"ca.crt":   tls.EncodeCertificatePEM(caCrt),
	}
	//Set owner reference
	addOwnerRefToObject(se, asOwner(br))
	err = action.Create(se)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	br.Spec.TLSSecretRef = &v1.ObjectReference{
		Name:      se.Name,
		Namespace: se.Namespace,
	}
	return nil
}

func newCACert(br *api.Broker) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tls.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	config := tls.CertConfig{
		CommonName:   fmt.Sprintf("%v operator CA", br.Name),
		Organization: []string{"starterpack.org"},
	}

	cert, err := tls.NewSelfSignedCACertificate(config, key)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, err
}

func newKeyAndCert(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, config tls.CertConfig) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := tls.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	cert, err := tls.NewSignedCertificate(config, key, caCert, caPrivKey)
	if err != nil {
		return nil, nil, err
	}
	return key, cert, nil
}
