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

# Assign a portal category to a V4 API, through the Terraform APIM provider.
# Categories are an inline attribute of apim_apiv4 (no standalone resource) and
# reference a category that must already exist in the environment (created as a
# precondition via mAPI in the scenario); APIM silently drops unknown refs.
# with_categories = false re-applies with an empty list to strip the reference.
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

variable "with_categories" {
  type    = bool
  default = true
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "categorized-api-tf"
  name            = "categorized-api-tf"
  description     = "V4 proxy API organised with a portal category"
  version         = "1"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"
  categories      = var.with_categories ? ["e2e-portal-category"] : []

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/categorized-api-tf/" }
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
  value = "/categorized-api-tf"
}
