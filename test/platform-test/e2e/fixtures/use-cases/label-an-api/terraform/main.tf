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

# Use case: organise a V4 API with labels, through the Terraform APIM provider.
# Labels are an inline attribute of apim_apiv4 (no standalone resource).
# with_labels = false re-applies with an empty list to strip them.
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

variable "with_labels" {
  type    = bool
  default = true
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-uc-label"
  name            = "e2e-tf-uc-label"
  description     = "E2E use-case: label a V4 API"
  version         = "1"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"
  labels          = var.with_labels ? ["e2e-uc-label-1", "e2e-uc-label-2"] : []

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/e2e-tf-uc-label/" }
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
      hrid     = "keyless"
      name     = "Free plan"
      type     = "API"
      mode     = "STANDARD"
      status   = "PUBLISHED"
      security = { type = "KEY_LESS" }
    }
  ]
}

output "api_id" {
  value = apim_apiv4.api.id
}

output "api_context_path" {
  value = "/e2e-tf-uc-label"
}
