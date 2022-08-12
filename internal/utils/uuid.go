package utils

import (
	"encoding/base64"

	uuid "github.com/satori/go.uuid" //nolint:gomodguard // to replace with google implementation
)

func ToUUID(decoded string) string {
	encoded := base64.RawStdEncoding.EncodeToString([]byte(decoded))
	return uuid.NewV3(uuid.NamespaceURL, encoded).String()
}

func NewUUID() string {
	return uuid.NewV4().String()
}
