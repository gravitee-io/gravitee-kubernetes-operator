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

terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

provider "apim" {
  # Configured via environment variables:
  #   APIM_SERVER_URL  (e.g. http://localhost:30083/automation)
  #   APIM_USERNAME
  #   APIM_PASSWORD
}

variable "environment_id" {
  type    = string
  default = "DEFAULT"
}

variable "organization_id" {
  type    = string
  default = "DEFAULT"
}

resource "apim_apiv4" "e2e_tf_gen_conditions" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-gen-conditions"
  name            = "e2e-tf-gen-conditions"
  description     = "E2E test: Terraform plan with general conditions page"
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
          {
            path = "/e2e-tf-gen-conditions/"
          }
        ]
        entrypoints = [
          {
            type = "http-proxy"
          }
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

  pages = [
    {
      hrid      = "e2e-gen-conds"
      name      = "E2E Plan General Conditions"
      type      = "MARKDOWN"
      content   = "These are the general conditions for the E2E test plan."
      published = true
    }
  ]

  plans = [
    {
      hrid                    = "keyless-with-conditions"
      name                    = "Keyless plan with conditions"
      description             = "A keyless plan that references a general conditions page"
      type                    = "API"
      mode                    = "STANDARD"
      validation              = "AUTO"
      status                  = "PUBLISHED"
      general_conditions_hrid = "e2e-gen-conds"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
}

output "api_id" {
  value = apim_apiv4.e2e_tf_gen_conditions.id
}
