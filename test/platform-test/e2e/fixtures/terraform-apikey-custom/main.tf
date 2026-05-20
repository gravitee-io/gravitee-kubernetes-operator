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

# GKO-2560: TF subscription on an API_KEY plan with a parameterised list of
# custom API keys. Driven from the test code via an auto.tfvars.json file so
# rotation/expireAt/length-boundary scenarios can re-apply the same workspace
# with different key sets.
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
  default = "custom"
}

# Parameterised key list. expire_at is optional per the TF provider schema
# (api_key_spec.go: ExpireAt is types.String, Optional). Defaulting to an
# empty list lets `terraform init`/`plan` run before the test writes vars.
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
  description     = "E2E test: TF API-key support (custom keys)"
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
  description     = "E2E test: TF API-key support — application"
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
  plan_hrid        = apim_apiv4.test.plans[0].hrid
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
