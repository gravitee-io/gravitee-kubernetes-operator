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

// +kubebuilder:object:generate=true
// +groupName=gravitee.io/v1alpha1

// Package docs holds the domain model for the Documentation CRD.
package docs

// The Go package is named "docs" rather than "documentation" (the CRD kind)
// because the go tool reserves the package name "documentation" (used to
// document non-Go programs) and silently ignores every file in such a package
// (see `go help packages`). Consumers alias this import back to "documentation".

// placeholder files for package level markers
