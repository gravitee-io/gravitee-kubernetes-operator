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

package uuid

import (
	"crypto/md5" //nolint:gosec // it is expected in this context
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/zeebo/xxh3"
)

func FromStrings(seeds ...string) string {
	seed := strings.Join(seeds, "")
	h := xxh3.HashString128(seed).Bytes()
	guid, err := uuid.FromBytes(h[:])
	if err != nil {
		panic(err)
	}
	return guid.String()
}

// JavaUUIDFromBytes generates a version 3 UUID based on the MD5 hash of the provided name, same as java.util.UUID.
func JavaUUIDFromBytes(data string) string {
	// Calculate the MD5 hash of the input name
	h := md5.New() //nolint:gosec // it is expected in this context
	h.Reset()
	h.Write([]byte(data))
	md5Bytes := h.Sum(nil)

	// Set version to 3 in the 7th byte (6th byte in zero-index).
	md5Bytes[6] &= 0x0f // clear the version bits
	md5Bytes[6] |= 0x30 // set to version 3

	// Set variant to IETF (RFC 4122) in the 9th byte (8th byte in zero-index).
	md5Bytes[8] &= 0x3f // clear the variant bits
	md5Bytes[8] |= 0x80 // set to IETF variant

	// Convert the byte array to a UUID string representation
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		md5Bytes[0:4],
		md5Bytes[4:6],
		md5Bytes[6:8],
		md5Bytes[8:10],
		md5Bytes[10:])
}

func NewV4String() string {
	return uuid.NewString()
}
