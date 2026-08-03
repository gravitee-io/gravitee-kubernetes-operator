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

# A V4 proxy API whose portal visibility and lifecycle state are driven by
# variables, so the journey can move it through the matrix with re-applies.
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

variable "visibility" {
  type    = string
  default = "PRIVATE"
}

variable "lifecycle_state" {
  type    = string
  default = "PUBLISHED"
}

# Gates the resource so the journey's remove("api") can drop it from the desired
# state and re-apply, the way a user deletes a resource block.
variable "create_api" {
  type    = bool
  default = true
}

resource "apim_apiv4" "api" {
  count           = var.create_api ? 1 : 0
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "portal-visibility-api-tf"
  name            = "portal-visibility-api-tf"
  description     = "V4 proxy API whose portal visibility and lifecycle state the journey changes"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = var.visibility
  lifecycle_state = var.lifecycle_state

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/portal-visibility-api-tf/" }
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
  value = var.create_api ? apim_apiv4.api[0].id : ""
}

output "api_context_path" {
  value = "/portal-visibility-api-tf"
}
