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

package apim

import (
	"github.com/gravitee-io/gravitee-kubernetes-operator/api/model/api/base"
	"github.com/gravitee-io/gravitee-kubernetes-operator/internal/apim/client"
)

// Notification brings support for notifications.
type Notification struct {
	*client.Client
}

func NewNotification(client *client.Client) *Notification {
	return &Notification{Client: client}
}

func (n *Notification) GetConsoleNotificationConfiguration(apiID string) (
	*base.ConsoleNotificationConfiguration, error) {
	notifications := make([]base.ConsoleNotificationConfiguration, 0)

	url := n.Client.EnvV1Target("apis").WithPath(apiID, "notificationsettings")
	err := n.Client.HTTP.Get(url.String(), &notifications)
	if err != nil {
		return nil, err
	}

	var console *base.ConsoleNotificationConfiguration
	for _, notification := range notifications {
		if notification.ConfigType == "PORTAL" {
			n := notification
			console = &n
			break
		}
	}
	return console, nil
}
