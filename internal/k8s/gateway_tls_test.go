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

package k8s

import (
	"context"
	"testing"

	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/gateway"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwAPIv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestGetServers_SingleTLSListener(t *testing.T) {
	gw := buildTestGateway(
		listener("https", gwAPIv1.HTTPSProtocolType, 443, tlsConfig("cert-a")),
	)

	portMapping := map[gwAPIv1.PortNumber]int32{443: 8443}
	servers, err := getServers(context.Background(), gw, portMapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(servers))
	}

	ssl, ok := servers[0]["ssl"].(map[string]any)
	if !ok {
		t.Fatal("expected ssl config on server")
	}

	if _, hasSNI := ssl["sni"]; hasSNI {
		t.Error("single TLS listener should NOT have sni enabled")
	}

	keystore := ssl["keystore"].(map[string]any)
	if keystore["secret"] != "secret://kubernetes/cert-a" {
		t.Errorf("expected secret://kubernetes/cert-a, got %v", keystore["secret"])
	}
}

func TestGetServers_MultipleTLSListenersSamePort(t *testing.T) {
	gw := buildTestGateway(
		listener("https-api", gwAPIv1.HTTPSProtocolType, 443, tlsConfig("cert-api")),
		listener("https-portal", gwAPIv1.HTTPSProtocolType, 443, tlsConfig("cert-portal")),
	)

	portMapping := map[gwAPIv1.PortNumber]int32{443: 8443}
	servers, err := getServers(context.Background(), gw, portMapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(servers) != 1 {
		t.Fatalf("expected 1 server (merged), got %d", len(servers))
	}

	ssl, ok := servers[0]["ssl"].(map[string]any)
	if !ok {
		t.Fatal("expected ssl config on server")
	}

	sni, hasSNI := ssl["sni"]
	if !hasSNI {
		t.Fatal("multiple TLS listeners on same port should have sni enabled")
	}
	if sni != true {
		t.Errorf("sni should be true, got %v", sni)
	}

	keystore := ssl["keystore"].(map[string]any)
	certs, hasCerts := keystore["certificates"].([]map[string]string)
	if !hasCerts {
		t.Fatal("expected certificates array in keystore")
	}
	if len(certs) != 2 {
		t.Fatalf("expected 2 certificates, got %d", len(certs))
	}
	if certs[0]["cert"] != "/opt/graviteeio-gateway/certs/https-api/tls.crt" {
		t.Errorf("unexpected cert path: %s", certs[0]["cert"])
	}
	if certs[1]["cert"] != "/opt/graviteeio-gateway/certs/https-portal/tls.crt" {
		t.Errorf("unexpected cert path: %s", certs[1]["cert"])
	}
}

func TestGetServers_TLSListenersDifferentPorts(t *testing.T) {
	gw := buildTestGateway(
		listener("https-api", gwAPIv1.HTTPSProtocolType, 443, tlsConfig("cert-api")),
		listener("https-admin", gwAPIv1.HTTPSProtocolType, 8443, tlsConfig("cert-admin")),
	)

	portMapping := map[gwAPIv1.PortNumber]int32{443: 8443, 8443: 9443}
	servers, err := getServers(context.Background(), gw, portMapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(servers) != 2 {
		t.Fatalf("expected 2 servers (different ports), got %d", len(servers))
	}

	for _, server := range servers {
		ssl := server["ssl"].(map[string]any)
		if _, hasSNI := ssl["sni"]; hasSNI {
			t.Error("separate ports should NOT have sni enabled")
		}
	}
}

func TestGetServers_MixedHTTPAndHTTPS(t *testing.T) {
	gw := buildTestGateway(
		listener("http", gwAPIv1.HTTPProtocolType, 80, nil),
		listener("https-api", gwAPIv1.HTTPSProtocolType, 443, tlsConfig("cert-api")),
		listener("https-portal", gwAPIv1.HTTPSProtocolType, 443, tlsConfig("cert-portal")),
	)

	portMapping := map[gwAPIv1.PortNumber]int32{80: 8080, 443: 8443}
	servers, err := getServers(context.Background(), gw, portMapping)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(servers) != 2 {
		t.Fatalf("expected 2 servers (HTTP + HTTPS), got %d", len(servers))
	}

	var httpServer, httpsServer map[string]any
	for _, s := range servers {
		if s["port"] == int32(8080) {
			httpServer = s
		} else {
			httpsServer = s
		}
	}

	if httpServer == nil {
		t.Fatal("expected HTTP server")
	}
	if _, hasSSL := httpServer["ssl"]; hasSSL {
		t.Error("HTTP server should not have ssl")
	}

	if httpsServer == nil {
		t.Fatal("expected HTTPS server")
	}
	ssl := httpsServer["ssl"].(map[string]any)
	if ssl["sni"] != true {
		t.Error("HTTPS server with multiple TLS listeners should have sni=true")
	}
}

func buildTestGateway(listeners ...gwAPIv1.Listener) *gateway.Gateway {
	gw := &gwAPIv1.Gateway{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "test-gw",
			Namespace: "default",
		},
		Spec: gwAPIv1.GatewaySpec{
			Listeners: listeners,
		},
		Status: gwAPIv1.GatewayStatus{
			Listeners: make([]gwAPIv1.ListenerStatus, len(listeners)),
		},
	}

	for i, l := range listeners {
		gw.Status.Listeners[i] = gwAPIv1.ListenerStatus{
			Name: l.Name,
			Conditions: []metaV1.Condition{
				{
					Type:   string(gwAPIv1.ListenerConditionAccepted),
					Status: metaV1.ConditionTrue,
				},
			},
		}
	}

	return gateway.WrapGateway(gw)
}

func listener(
	name string,
	protocol gwAPIv1.ProtocolType,
	port gwAPIv1.PortNumber,
	tls *gwAPIv1.ListenerTLSConfig,
) gwAPIv1.Listener {
	return gwAPIv1.Listener{
		Name:     gwAPIv1.SectionName(name),
		Protocol: protocol,
		Port:     port,
		TLS:      tls,
	}
}

func tlsConfig(secretName string) *gwAPIv1.ListenerTLSConfig {
	return &gwAPIv1.ListenerTLSConfig{
		CertificateRefs: []gwAPIv1.SecretObjectReference{
			{Name: gwAPIv1.ObjectName(secretName)},
		},
	}
}
