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

# Publish a V4 proxy API to the portal and serve traffic, via the Terraform APIM
# provider. `state` is re-applied (STARTED <-> STOPPED) to take the API on and off
# the gateway.
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

variable "state" {
  type    = string
  default = "STARTED"
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "published-api-tf"
  name            = "published-api-tf"
  description     = "V4 proxy API published to the portal, verifying start/stop and gateway reachability"
  version         = "1"
  type            = "PROXY"
  state           = var.state
  lifecycle_state = "PUBLISHED"
  visibility      = "PUBLIC"

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/published-api-tf/" }
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
          name                  = "default-endpoint"
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

output "api_context_path" {
  value = "/published-api-tf"
}
