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

# GKO-2560: TF API with two coexisting plans (KEY_LESS + API_KEY). The
# subscription targets the api-key plan with a parameterised key list. Mirrors
# the GKO `V4_APIKEY_MIXED_WITH_KEYLESS` scenario so the gateway routing
# discriminator can be tested through the TF write path.
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
  default = "mixed"
}

variable "keys" {
  type = list(object({
    key       = string
    expire_at = optional(string)
  }))
  default = []
}

resource "apim_apiv4" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-apikey-${var.hrid_suffix}"
  name            = "e2e-tf-apikey-${var.hrid_suffix}"
  description     = "E2E test: TF API-key + keyless plan coexistence"
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
          { path = "/e2e-tf-apikey-${var.hrid_suffix}/" }
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
      name       = "Keyless plan"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      security = {
        type = "KEY_LESS"
      }
    },
    {
      hrid       = "apikey"
      name       = "Api Key plan"
      type       = "API"
      mode       = "STANDARD"
      validation = "AUTO"
      status     = "PUBLISHED"
      security = {
        type = "API_KEY"
      }
    }
  ]
}

resource "apim_application" "test" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "e2e-tf-apikey-app-${var.hrid_suffix}"
  name            = "e2e-tf-apikey-app-${var.hrid_suffix}"
  description     = "E2E test: TF API-key + keyless plan — application"
  settings = {
    app = {
      type = "SIMPLE"
    }
  }
}

resource "apim_subscription" "test" {
  environment_id   = var.environment_id
  organization_id  = var.organization_id
  hrid             = "e2e-tf-apikey-sub-${var.hrid_suffix}"
  api_hrid         = apim_apiv4.test.hrid
  plan_hrid        = "apikey"
  application_hrid = apim_application.test.hrid
  api_keys         = var.keys
}

output "api_id" {
  value = apim_apiv4.test.id
}

output "sub_id" {
  value = apim_subscription.test.id
}

output "api_context_path" {
  value = "/e2e-tf-apikey-${var.hrid_suffix}"
}
