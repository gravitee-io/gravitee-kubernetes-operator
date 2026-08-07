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

# An API with a JWT plan and an OAuth2 plan (both manually validated), an
# application, and a subscription to each plan. Both plans are MANUAL on purpose:
# a subscription written through the Automation API is auto-validated anyway,
# which is what the journey asserts.
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

resource "apim_apiv4" "api" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "secured-plans-api-tf"
  name            = "secured-plans-api-tf"
  description     = "V4 proxy API exposing a JWT plan and an OAuth2 plan, both manually validated"
  version         = "1.0"
  type            = "PROXY"
  state           = "STARTED"
  lifecycle_state = "PUBLISHED"
  visibility      = "PRIVATE"

  # APIM otherwise rejects the second subscription with
  # "An other OAuth2 or JWT plan is already subscribed by the same application".
  allow_multi_jwt_oauth2_subscriptions = true
  # Declared alongside it because both govern how this API may be consumed. APIM
  # exposes no product surface through either driver, so the journey can only
  # assert that the flag round-trips.
  allowed_in_api_products = true

  listeners = [
    {
      http = {
        type = "HTTP"
        paths = [
          { path = "/secured-plans-api-tf/" }
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
      hrid       = "jwt-plan"
      name       = "JWT plan"
      type       = "API"
      mode       = "STANDARD"
      validation = "MANUAL"
      status     = "PUBLISHED"
      security   = { type = "JWT" }
    },
    {
      hrid       = "oauth2-plan"
      name       = "OAuth2 plan"
      type       = "API"
      mode       = "STANDARD"
      validation = "MANUAL"
      status     = "PUBLISHED"
      security   = { type = "OAUTH2" }
    }
  ]
}

resource "apim_application" "app" {
  environment_id  = var.environment_id
  organization_id = var.organization_id
  hrid            = "secured-plans-app-tf"
  name            = "secured-plans-app-tf"
  description     = "Consumer application subscribing to the journey's secured plans"
  settings = {
    app = {
      type      = "SIMPLE"
      client_id = "secured-plans-app-tf-client"
    }
  }
}

resource "apim_subscription" "jwt" {
  environment_id   = var.environment_id
  organization_id  = var.organization_id
  hrid             = "secured-plans-sub-jwt-tf"
  api_hrid         = apim_apiv4.api.hrid
  application_hrid = apim_application.app.hrid
  plan_hrid        = "jwt-plan"
}

resource "apim_subscription" "oauth2" {
  environment_id   = var.environment_id
  organization_id  = var.organization_id
  hrid             = "secured-plans-sub-oauth2-tf"
  api_hrid         = apim_apiv4.api.hrid
  application_hrid = apim_application.app.hrid
  plan_hrid        = "oauth2-plan"
}

output "api_id" {
  value = apim_apiv4.api.id
}

output "app_id" {
  value = apim_application.app.id
}

output "sub_jwt_id" {
  value = apim_subscription.jwt.id
}

output "sub_oauth2_id" {
  value = apim_subscription.oauth2.id
}

output "api_context_path" {
  value = "/secured-plans-api-tf"
}
