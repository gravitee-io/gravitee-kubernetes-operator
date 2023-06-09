// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"bufio"
	"bytes"
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	kerrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/gateway"
	"github.com/gravitee-io/gravitee-kubernetes-operator/pkg/keys"
	netV1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ks "github.com/pavlo-v-chernykh/keystore-go/v4"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

const graviteeConfigFile = "gravitee.yml"

type keystoreCredentials struct {
	name string
	key  string
	pass []byte
}

func (d *Delegate) retrieveIngressListWithTLS(ctx context.Context, ns string) (*netV1.IngressList, error) {
	il := &netV1.IngressList{}
	if err := d.k8s.List(ctx, il, client.InNamespace(ns)); err != nil {
		return nil, client.IgnoreNotFound(err)
	}

	result := &netV1.IngressList{}
	for i := range il.Items {
		ingress := il.Items[i]
		if ingress.GetAnnotations()[keys.IngressClassAnnotation] == keys.IngressClassAnnotationValue {
			if ingress.Spec.TLS != nil {
				result.Items = append(result.Items, ingress)
			}
		}
	}

	return result, nil
}

func (d *Delegate) removeKeyFromKeystore(secret *v1.Secret) error {
	ksc, err := d.getKeystoreCredentials(secret.Namespace)
	if err != nil {
		return err
	}

	pki, _, err := d.generateKeyPair(secret)
	if err != nil {
		return err
	}

	nn := &types.NamespacedName{Namespace: secret.Namespace, Name: ksc.name}
	gwKeyStoreSecret, jks, err := d.readKeyStore(nn, ksc)
	if err != nil {
		return err
	}

	jks.DeleteEntry(pki.CommonName)

	return d.writeToKeyStore(gwKeyStoreSecret, jks, ksc)
}

func (d *Delegate) updateKeyInKeystore(secret *v1.Secret) error {
	ksc, err := d.getKeystoreCredentials(secret.Namespace)
	if err != nil {
		return err
	}

	pki, keyPair, err := d.generateKeyPair(secret)
	if err != nil {
		return err
	}

	nn := &types.NamespacedName{Namespace: secret.Namespace, Name: ksc.name}
	gwKeyStoreSecret, jks, err := d.readKeyStore(nn, ksc)
	if err != nil {
		return err
	}

	if err = jks.SetPrivateKeyEntry(pki.CommonName, *keyPair, ksc.pass); err != nil {
		return err
	}

	return d.writeToKeyStore(gwKeyStoreSecret, jks, ksc)
}

// returns the name of gw keystore and the password to open it.
func (d *Delegate) getKeystoreCredentials(ns string) (*keystoreCredentials, error) {
	// This secret will give us the name and the password for opening the gateway keystore
	// The keystore should be jks format
	if ksc, err := d.autoDiscoverGatewayKeystore(ns); client.IgnoreNotFound(err) != nil {
		return nil, err
	} else if ksc != nil {
		return ksc, nil
	}

	sl := &v1.SecretList{}
	if err := d.k8s.List(
		d.ctx, sl,
		client.InNamespace(ns),
		client.MatchingLabels{keys.GatewayKeystoreConfigSecret: "true"}); err != nil {
		return nil, client.IgnoreNotFound(err)
	}

	if len(sl.Items) == 0 {
		return nil, fmt.Errorf("%s %s %s:%s", "can't find a secret for accessing the gateway keystore",
			"you need to label you secret with", keys.GatewayKeystoreConfigSecret, "true")
	} else if len(sl.Items) > 1 {
		return nil, fmt.Errorf("%s %s", "found more than one secrets with label", keys.GatewayKeystoreConfigSecret)
	}

	s := sl.Items[0]
	if len(s.Data) == 0 {
		return nil, fmt.Errorf("no credentials provided to access the gateway keystore")
	}

	return &keystoreCredentials{name: string(s.Data["name"]), key: string(s.Data["key"]), pass: s.Data["password"]}, nil
}

func (d *Delegate) autoDiscoverGatewayKeystore(ns string) (*keystoreCredentials, error) {
	// get gravitee.yml from the configmap
	cfg := &gateway.Config{}
	if err := d.unmarshalGatewayConfig(ns, cfg); err != nil {
		return nil, err
	}

	ks := cfg.HTTP.TLS.Keystore

	if err := ks.Validate(); err != nil {
		return nil, err
	}

	ksNS := ks.Location.Namespace()
	ksName := ks.Location.Name()
	ksKey := ks.Location.Key()
	ksPassword := ks.Password

	if ksNS != ns {
		return nil, fmt.Errorf("keystore is outside of the current namespace")
	}

	kubernetesPassword := gateway.GraviteeKubeProperty(ksPassword)

	if !kubernetesPassword.IsValid() {
		return &keystoreCredentials{
			name: ksName,
			key:  ksKey,
			pass: []byte(ksPassword),
		}, nil
	}

	if kubernetesPassword.Namespace() != ns {
		return nil, fmt.Errorf(
			"password location is outside of the current namespace",
		)
	}

	password, err := d.resolveKubernetesPassword(kubernetesPassword)
	if err != nil {
		return nil, err
	}

	return &keystoreCredentials{
		name: ksName,
		key:  ksKey,
		pass: password,
	}, nil
}

func (d *Delegate) resolveKubernetesPassword(prop gateway.GraviteeKubeProperty) ([]byte, error) {
	location := types.NamespacedName{
		Namespace: prop.Namespace(),
		Name:      prop.Name(),
	}
	obj := prop.NewReceiver()
	if err := d.k8s.Get(d.ctx, location, obj); err != nil {
		return nil, err
	}
	password := prop.Get(obj)
	if password == nil {
		return nil, fmt.Errorf("can't resolve password from %s", location)
	}
	return password, nil
}

func (d *Delegate) unmarshalGatewayConfig(ns string, cfg *gateway.Config) error {
	cl := &v1.ConfigMapList{}
	if err := d.k8s.List(
		d.ctx, cl,
		client.InNamespace(ns),
		client.MatchingLabels{
			keys.GraviteeComponentLabel: keys.IngressComponentLabelValue,
		}); err != nil {
		return client.IgnoreNotFound(err)
	}

	if len(cl.Items) != 1 || cl.Items[0].Data[graviteeConfigFile] == "" {
		d.log.Info("can't automatically find gateway gravitee.yml config")
		return kerrors.NewNotFound(v1.Resource(graviteeConfigFile), graviteeConfigFile)
	}

	return yaml.Unmarshal([]byte(cl.Items[0].Data[graviteeConfigFile]), cfg)
}

// convert K8S tls secret to a keypair.
func (d *Delegate) generateKeyPair(secret *v1.Secret) (*pkix.Name, *ks.PrivateKeyEntry, error) {
	// get the key and certificate (The TLS secret must contain keys named tls.crt and tls.key
	// https://kubernetes.io/docs/concepts/services-networking/ingress/#tls
	pemKeyBytes, ok := secret.Data["tls.key"]
	if !ok {
		return nil, nil, fmt.Errorf("%s", "tls key not found in secret")
	}

	tlsKey, _ := pem.Decode(pemKeyBytes)
	if tlsKey == nil {
		return nil, nil, fmt.Errorf("%s", "can not decode the tls key")
	}

	if pkcs1, err := x509.ParsePKCS1PrivateKey(tlsKey.Bytes); err == nil {
		pkcs8, mErr := x509.MarshalPKCS8PrivateKey(pkcs1)
		if mErr != nil {
			return nil, nil, mErr
		}
		tlsKey.Bytes = pkcs8
	}

	if !strings.Contains(tlsKey.Type, "PRIVATE KEY") {
		return nil, nil, fmt.Errorf("%s", "wrong tls key type")
	}

	pemCrtBytes, ok := secret.Data["tls.crt"]
	if !ok {
		return nil, nil, fmt.Errorf("%s", "tls cert not found in secret")
	}

	tlsCrt, _ := pem.Decode(pemCrtBytes)
	if tlsCrt == nil {
		return nil, nil, fmt.Errorf("%s", "can not decode the tls certificate")
	}

	cert, err := x509.ParseCertificate(tlsCrt.Bytes)
	if err != nil {
		return nil, nil, err
	}

	if tlsCrt.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("%s", "wrong tls certification type")
	}

	pke := &ks.PrivateKeyEntry{
		CreationTime: time.Now(),
		PrivateKey:   tlsKey.Bytes,
		CertificateChain: []ks.Certificate{
			{
				Type:    "X509",
				Content: tlsCrt.Bytes,
			},
		},
	}

	return &cert.Subject, pke, nil
}

func (d *Delegate) readKeyStore(nn *types.NamespacedName, ksc *keystoreCredentials) (*v1.Secret, *ks.KeyStore, error) {
	gwKeystoreSecret := &v1.Secret{}
	if err := d.k8s.Get(d.ctx, *nn, gwKeystoreSecret); err != nil {
		return nil, nil, err
	}

	data := gwKeystoreSecret.Data[ksc.key]
	if data == nil {
		return nil, nil, fmt.Errorf("unable to find keystore data for the gateway")
	}

	jks := ks.New()
	if err := jks.Load(bytes.NewReader(data), ksc.pass); err != nil { // should come from a variable
		return nil, nil, err
	}

	return gwKeystoreSecret, &jks, nil
}

func (d *Delegate) writeToKeyStore(ksSecret *v1.Secret, jks *ks.KeyStore, ksc *keystoreCredentials) error {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err := jks.Store(writer, ksc.pass)
	if err != nil {
		return err
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	ksSecret.Data[ksc.key] = b.Bytes()
	return d.k8s.Update(d.ctx, ksSecret)
}
