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

# A DYNAMIC dictionary whose HTTP provider polls the Gravitee echo endpoint every
# 5s; a JOLT spec maps the echoed X-Test-Specific header into a property of the
# same name. A keyless PROXY API injects that property into the X-Env response
# header via transform-headers. The lifecycle knobs are driven by tfvars written
# by the test harness (Provisioned.update / .remove):
#   - header_value      : the provider request-header value (update propagation)
#   - deployed          : deploy/undeploy the dictionary (deployed=false stop)
#   - create_dictionary : count-gate the dictionary (delete-stops via remove())
#
# The API's EL key is the dictionary HRID as a LITERAL string, not a reference to
# apim_dictionary.dyn[0].hrid: a reference would make the API depend on the
# count-gated dictionary, so dropping the dictionary (create_dictionary=false)
# would force the API to be replaced instead of just stopping resolution.
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

variable "header_value" {
  type    = string
  default = "ABCDEF"
}

variable "deployed" {
  type    = bool
  default = true
}

variable "create_dictionary" {
  type    = bool
  default = true
}

resource "apim_dictionary" "dyn" {
  count           = var.create_dictionary ? 1 : 0
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "dyn-dictionary-tf"
  name            = "dyn-dictionary-tf"
  description     = "DYNAMIC dictionary exposing echo headers as properties"
  type            = "DYNAMIC"
  deployed        = var.deployed
  dynamic = {
    provider = {
      http = {
        type   = "HTTP"
        url    = "https://api.gravitee.io/echo"
        method = "GET"
        headers = [
          {
            name  = "X-Test-Specific"
            value = var.header_value
          }
        ]
        specification = <<-EOT
        [
          {
            "operation": "shift",
            "spec": {
              "headers": {
                "*": {
                  "$": "[#2].key",
                  "@": "[#2].value"
                }
              }
            }
          }
        ]
        EOT
      }
    }
    trigger = {
      rate = 5
      unit = "SECONDS"
    }
  }
}

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "dyn-dictionary-consumer-api-tf"
  name            = "dyn-dictionary-consumer-api-tf"
  description     = "API injecting a dynamic dictionary value via transform-headers"
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
          { path = "/dyn-dictionary-consumer-api-tf" }
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
      name    = "Add dynamic dictionary header"
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
                name  = "X-Env"
                value = "{#dictionaries['dyn-dictionary-tf']['X-Test-Specific']}"
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
  value = "/dyn-dictionary-consumer-api-tf"
}
