# Copyright (C) 2015 The Gravitee team (http://gravitee.io)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# A group and a V4 proxy API whose console notification targets it. Terraform
# inlines what GKO models as a standalone Notification CR, so the event list is a
# variable. Referring to the group through `apim_group.group.name` also gives
# Terraform the dependency edge that guarantees the group exists before the
# notification targets it.
terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {}

variable "environment_id" {
  type    = string
  default = "DEFAULT"
}

variable "organization_id" {
  type    = string
  default = "DEFAULT"
}

variable "notification_events" {
  type    = list(string)
  default = ["API_STARTED", "API_STOPPED"]
}

resource "apim_group" "group" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "notified-group-tf"
  name            = "notified-group-tf"
  notify_members  = false
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "notified-api-tf"
  name            = "notified-api-tf"
  description     = "V4 proxy API whose console notification settings the journey configures"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"
  notify_members  = false

  # APIM only allows a console notification to target groups the API itself is
  # associated with, so the association is part of the setup, not the assertion.
  groups = [apim_group.group.name]

  console_notification = {
    events = var.notification_events
    groups = [apim_group.group.name]
  }

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/notified-api-tf/" }
        ]
        entrypoints = [
          { type = "http-proxy" }
        ]
      }
    }
  ]

  endpoint_groups = [
    {
      name = "Default HTTP proxy group"
      type = "http-proxy"
      endpoints = [
        {
          name                  = "Default HTTP proxy"
          type                  = "http-proxy"
          inherit_configuration = false
          configuration         = jsonencode({ target = "https://api.gravitee.io/echo" })
        }
      ]
    }
  ]

  flow_execution = {
    mode           = "DEFAULT"
    match_required = false
  }

  plans = [
    {
      hrid       = "keyless"
      name       = "Free plan"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}

output "api_id" {
  value = apim_apiv4.api.id
}

output "group_id" {
  value = apim_group.group.id
}

output "api_context_path" {
  value = "/notified-api-tf"
}
