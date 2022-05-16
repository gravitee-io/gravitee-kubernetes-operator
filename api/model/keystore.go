package model

type KeyStoreType int

const (
	PEM KeyStoreType = iota
	PKCS12
	JKS
)

type TrustStore struct {
	TrustStoreType KeyStoreType `json:"type,omitempty"`
}

type KeyStore struct {
	KeyStoreType KeyStoreType `json:"type,omitempty"`
}
type PEMTrustStore struct {
	Type    KeyStoreType `json:"type,omitempty"`
	Path    string       `json:"path,omitempty"`
	Content string       `json:"content,omitempty"`
}

type PKCS12TrustStore struct {
	Type     KeyStoreType `json:"type,omitempty"`
	Path     string       `json:"path,omitempty"`
	Content  string       `json:"content,omitempty"`
	Password string       `json:"password,omitempty"`
}
type JKSTrustStore struct {
	Type     KeyStoreType `json:"type,omitempty"`
	Path     string       `json:"path,omitempty"`
	Content  string       `json:"content,omitempty"`
	Password string       `json:"password,omitempty"`
}
type PEMKeyStore struct {
	Type        KeyStoreType `json:"type,omitempty"`
	KeyPath     string       `json:"keyPath,omitempty"`
	KeyContent  string       `json:"keyContent,omitempty"`
	CertPath    string       `json:"certPath,omitempty"`
	CertContent string       `json:"certContent,omitempty"`
}

type PKCS12KeyStore struct {
	Type     KeyStoreType `json:"type,omitempty"`
	Path     string       `json:"path,omitempty"`
	Content  string       `json:"content,omitempty"`
	Password string       `json:"password,omitempty"`
}

type JKSKeyStore struct {
	Type     KeyStoreType `json:"type,omitempty"`
	Path     string       `json:"path,omitempty"`
	Content  string       `json:"content,omitempty"`
	Password string       `json:"password,omitempty"`
}
