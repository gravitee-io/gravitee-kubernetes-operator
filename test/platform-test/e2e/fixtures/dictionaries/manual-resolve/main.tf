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

# Terraform arm of the cross-provisioner `dictionaries/manual-resolve` scenario:
# a MANUAL dictionary plus a keyless API whose transform-headers flow injects the
# dictionary value into a response header. The shared scenario body asserts the
# gateway resolves `{#dictionaries['<hrid>']['env']}` to "test". Mirrors the GKO
# arm (fixtures/dictionaries/{dictionary-manual,api-with-dictionary}/crd.yaml).
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

variable "hrid_suffix" {
  type    = string
  default = "resolve"
}

resource "apim_dictionary" "dict" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-dict-${var.hrid_suffix}"
  name            = "e2e-tf-dict-${var.hrid_suffix}"
  description     = "E2E test: TF MANUAL dictionary resolved at the gateway"
  type            = "MANUAL"
  deployed        = true
  manual = {
    properties = {
      env = "test"
    }
  }
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-dictapi-${var.hrid_suffix}"
  name            = "e2e-tf-dictapi-${var.hrid_suffix}"
  description     = "E2E test: API injecting a dictionary value via transform-headers"
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
          { path = "/e2e-tf-dictapi-${var.hrid_suffix}" }
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

  # Inject the dictionary value into a request header the echo endpoint reflects.
  # The EL key is the dictionary HRID (APIM keys deployed dictionaries by HRID).
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
  value = "/e2e-tf-dictapi-${var.hrid_suffix}"
}
