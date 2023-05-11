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
	keystoreYaml, err := d.unmarshalGatewayConfig(ns)
	if err != nil {
		return nil, err
	}

	ksType, ok := keystoreYaml["type"].(string)
	if !ok || ksType == "" {
		return nil, fmt.Errorf("%s doesn't include a http.ssl.keystore.type", graviteeConfigFile)
	}

	if ksType != "jks" {
		return nil, fmt.Errorf("unsupported keystore type. GKO only supports jks keystores at this moment")
	}

	kubernetes, ok := keystoreYaml["kubernetes"].(string)
	if !ok || kubernetes == "" {
		return nil, fmt.Errorf("%s doesn't include a http.ssl.keystore.kubernetes", graviteeConfigFile)
	}

	password, ok := keystoreYaml["password"].(string)
	if !ok || password == "" {
		return nil, fmt.Errorf("%s doesn't include a http.ssl.keystore.password", graviteeConfigFile)
	}

	ksName := strings.Split(kubernetes, "/")
	// example: /default/secrets/api-custom-cert-opaque/keystore
	k8sPropertyLength := 5
	if len(ksName) != k8sPropertyLength {
		return nil, fmt.Errorf("wrong keystore name. it should be like /${NAMESPACE}/secrets/${SECRET_NAME}/${KEY_NAME}")
	}

	if ksName[1] != ns {
		return nil, fmt.Errorf("keystore is outside of the current namespace")
	}

	if strings.HasPrefix(password, "kubernetes://") {
		return d.resolveKubernetesProperty(ns, password, ksName)
	}

	return &keystoreCredentials{name: ksName[3], key: ksName[4], pass: []byte(password)}, nil
}

func (d *Delegate) resolveKubernetesProperty(ns string, pass string, ksName []string) (*keystoreCredentials, error) {
	// kubernetes properties must have a length of 6 if you split them with "kubernetes://"
	// example: kubernetes://default/secret/gateway-secret/my_key
	ksPass := strings.Split(pass, "/")
	ksPropertyLength := 6
	if len(ksPass) != ksPropertyLength {
		return nil, fmt.Errorf("%s%s", "wrong Gateway keystore password",
			"if you reference a secret it should be like	kubernetes://${NAMESPACE}/secrets/${SECRET_NAME}/${SECRET_KEY}")
	}

	if ksPass[0] != "kubernetes:" {
		return nil, fmt.Errorf("unsupported property type %s", ksPass[0])
	}

	if ksPass[2] != ns {
		return nil, fmt.Errorf("keystore password is outside the current namespace")
	}

	switch ksPass[3] {
	case "secret":
		sec := new(v1.Secret)
		if err := d.k8s.Get(d.ctx, types.NamespacedName{Namespace: ns, Name: ksPass[4]}, sec); err != nil {
			return nil, err
		}
		return &keystoreCredentials{name: ksName[3], key: ksName[4], pass: sec.Data[ksPass[5]]}, nil
	case "configmaps":
		cm := new(v1.ConfigMap)
		if err := d.k8s.Get(d.ctx, types.NamespacedName{Namespace: ns, Name: ksPass[4]}, cm); err != nil {
			return nil, err
		}
		return &keystoreCredentials{name: ksName[3], key: ksName[4], pass: []byte(cm.Data[ksPass[5]])}, nil
	default:
		return nil, fmt.Errorf("unsupported resource type %s", ksPass[3])
	}
}

func (d *Delegate) unmarshalGatewayConfig(ns string) (map[string]interface{}, error) {
	cl := &v1.ConfigMapList{}
	if err := d.k8s.List(
		d.ctx, cl,
		client.InNamespace(ns),
		client.MatchingLabels{
			keys.GraviteeComponentLabel: keys.IngressComponentLabelValue,
		}); err != nil {
		return nil, client.IgnoreNotFound(err)
	}

	if len(cl.Items) != 1 || cl.Items[0].Data[graviteeConfigFile] == "" {
		d.log.Info("can't automatically find gateway gravitee.yml config")
		return nil, kerrors.NewNotFound(v1.Resource(graviteeConfigFile), graviteeConfigFile)
	}

	yml := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(cl.Items[0].Data[graviteeConfigFile]), yml); err != nil {
		return nil, err
	}

	http, ok := yml["http"].(map[string]interface{})
	if !ok || http == nil {
		return nil, fmt.Errorf("%s doesn't include a http section", graviteeConfigFile)
	}

	ssl, ok := http["ssl"].(map[string]interface{})
	if !ok || ssl == nil {
		return nil, fmt.Errorf("%s doesn't include a http.ssl section", graviteeConfigFile)
	}

	ks, ok := ssl["keystore"].(map[string]interface{})
	if !ok || ks == nil {
		return nil, fmt.Errorf("%s doesn't include a http.ssl.keystore section", graviteeConfigFile)
	}

	return ks, nil
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
	err = writer.Flush()
	if err != nil {
		return err
	}

	ksSecret.Data[ksc.key] = b.Bytes()

	return d.k8s.Update(d.ctx, ksSecret)
}
