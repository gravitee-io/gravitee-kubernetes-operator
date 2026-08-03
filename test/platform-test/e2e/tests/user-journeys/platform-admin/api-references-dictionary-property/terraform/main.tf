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

# A MANUAL dictionary plus a keyless API whose transform-headers flow injects the
# dictionary value into a response header, via the Terraform APIM provider. The
# EL key is the dictionary HRID (APIM keys deployed dictionaries by HRID).
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

# Value of the dictionary's `env` property. Defaults to "test" so the resolve
# scenario (no tfvars) is unchanged; the update scenario overrides it via tfvars.
variable "env_value" {
  type    = string
  default = "test"
}

resource "apim_dictionary" "dict" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "manual-dictionary-tf"
  name            = "manual-dictionary-tf"
  description     = "MANUAL dictionary resolved at the gateway"
  type            = "MANUAL"
  deployed        = true
  manual = {
    properties = {
      env = var.env_value
    }
  }
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "dictionary-consumer-api-tf"
  name            = "dictionary-consumer-api-tf"
  description     = "API injecting a dictionary value via transform-headers"
  version         = "1"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/dictionary-consumer-api-tf" }
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

  flows = [
    {
      enabled = true
      name    = "Add dictionary header"
      selectors = [
        {
          http = {
            type          = "HTTP"
            path          = "/"
            path_operator = "STARTS_WITH"
          }
        }
      ]
      request = [
        {
          enabled = true
          name    = "Transform Headers"
          policy  = "transform-headers"
          configuration = jsonencode({
            addHeaders = [
              {
                name  = "X-Dict-Env"
                value = "{#dictionaries['${apim_dictionary.dict.hrid}']['env']}"
              }
            ]
          })
        }
      ]
    }
  ]

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
  value = "/dictionary-consumer-api-tf"
}
