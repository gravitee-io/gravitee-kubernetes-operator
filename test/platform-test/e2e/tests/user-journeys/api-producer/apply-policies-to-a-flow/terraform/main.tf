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

# A V4 proxy API with a transform-headers policy on the response phase. The flow
# name and the header it adds are variables so the journey can rewrite the policy
# and drop it entirely with re-applies. Flows are inline on apim_apiv4.
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

variable "with_policy" {
  type    = bool
  default = true
}

variable "flow_name" {
  type    = string
  default = "Add custom header"
}

variable "header_name" {
  type    = string
  default = "X-E2E-Policy"
}

variable "header_value" {
  type    = string
  default = "applied"
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "policy-api-tf"
  name            = "policy-api-tf"
  description     = "V4 proxy API whose response-phase flow the journey adds, changes and removes"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"

  flows = var.with_policy ? [
    {
      name    = var.flow_name
      enabled = true
      selectors = [
        {
          http = {
            type          = "HTTP"
            path          = "/"
            path_operator = "STARTS_WITH"
          }
        }
      ]
      response = [
        {
          name    = "Transform Headers"
          enabled = true
          policy  = "transform-headers"
          configuration = jsonencode({
            addHeaders = [
              { name = var.header_name, value = var.header_value }
            ]
          })
        }
      ]
    }
  ] : []

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/policy-api-tf/" }
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

output "api_context_path" {
  value = "/policy-api-tf"
}
