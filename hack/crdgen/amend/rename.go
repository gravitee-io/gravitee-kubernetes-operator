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

package amend

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// renameKey renames a component schema. x-go-type-name cannot do it: oapi-codegen keeps the OAS name for the
// fields and adds the new one as an alias. x-go-name cannot either: it also renames the fields holding the schema.
const renameKey = "x-gko-name"

const schemaRefPrefix = "#/components/schemas/"

// RenameSchemas renames each component schema carrying x-gko-name, rewrites every $ref to it and drops the key.
// It runs on the YAML tree, before the spec is loaded, so that every $ref is still a plain string.
func RenameSchemas(root *yaml.Node) error {
	schemas := mappingValue(mappingValue(document(root), "components"), "schemas")
	if schemas == nil {
		return nil
	}

	refs := map[string]string{}
	for i := 0; i+1 < len(schemas.Content); i += 2 {
		key, schema := schemas.Content[i], schemas.Content[i+1]
		name, err := popScalar(schema, renameKey)
		if err != nil {
			return fmt.Errorf("%s: %s: %w", key.Value, renameKey, err)
		}
		if name == "" {
			continue
		}
		refs[schemaRefPrefix+key.Value] = schemaRefPrefix + name
		key.Value = name
	}

	seen := map[string]bool{}
	for i := 0; i < len(schemas.Content); i += 2 {
		name := schemas.Content[i].Value
		if seen[name] {
			return fmt.Errorf("%s: two schemas have this name after %s", name, renameKey)
		}
		seen[name] = true
	}

	rewriteRefs(root, refs)
	return nil
}

func document(root *yaml.Node) *yaml.Node {
	if root != nil && root.Kind == yaml.DocumentNode && len(root.Content) > 0 {
		return root.Content[0]
	}
	return root
}

// mappingValue returns the value of key in a mapping node, or nil.
func mappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

// popScalar removes key from a mapping node and returns its string value, "" when absent.
func popScalar(node *yaml.Node, key string) (string, error) {
	if node.Kind != yaml.MappingNode {
		return "", nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value != key {
			continue
		}
		value := node.Content[i+1]
		if value.Kind != yaml.ScalarNode || value.Value == "" {
			return "", fmt.Errorf("want a schema name, got %v", value.Value)
		}
		node.Content = append(node.Content[:i], node.Content[i+2:]...)
		return value.Value, nil
	}
	return "", nil
}

func rewriteRefs(node *yaml.Node, refs map[string]string) {
	if node.Kind == yaml.MappingNode {
		for i := 0; i+1 < len(node.Content); i += 2 {
			key, value := node.Content[i], node.Content[i+1]
			if renamed, ok := refs[value.Value]; ok && key.Value == "$ref" && value.Kind == yaml.ScalarNode {
				value.Value = renamed
			}
		}
	}
	for _, child := range node.Content {
		rewriteRefs(child, refs)
	}
}
