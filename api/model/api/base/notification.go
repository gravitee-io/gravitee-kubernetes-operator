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

package base

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/notification"
)

// NotificationConfigurationBase base object for notifications.
type NotificationConfigurationBase struct {
	ReferenceType string   `json:"referenceType,omitempty"`
	ReferenceID   string   `json:"referenceId,omitempty"`
	Hooks         []string `json:"hooks"`
	ConfigType    string   `json:"config_type"`
	Origin        string   `json:"origin"`
}

// ConsoleNotificationConfiguration mAPI object to update notification settings.
type ConsoleNotificationConfiguration struct {
	NotificationConfigurationBase `json:",inline"`
	User                          string   `json:"user,omitempty"`
	Groups                        []string `json:"groups"`
}

func (s *ConsoleNotificationConfiguration) IsConsoleNotification() bool {
	// hack to have an interface
	return true
}

// ToAPIConsoleNotificationSettings transforms a Console object into a ConsoleNotificationConfiguration object for API usage.
// It set the configType to "PORTAL".
func ToAPIConsoleNotificationSettings(console notification.Console) *ConsoleNotificationConfiguration {
	var configuration ConsoleNotificationConfiguration
	configuration.Hooks = console.APIEventsAsString()
	configuration.Groups = console.Groups
	configuration.ConfigType = "PORTAL"
	configuration.Origin = "KUBERNETES"
	return &configuration
}
